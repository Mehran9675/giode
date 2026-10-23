// Package scroll provides a stateful vertical scroll container.
package scroll

import (
	"image"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
)

const (
	// unbounded mirrors the sentinel Gio uses for infinite
	// constraints.
	unbounded = 1e6
)

// Scroll is a stateful vertical scroll container. It implements
// elements.Element, so it can be placed directly in the tree like any
// other element. Create it once with New and lay it out every frame;
// content is rebuilt every frame while visible and may be taller than
// the viewport.
//
// The scroll offset reacts to the mouse wheel and touch drags, and
// the scrollbar thumb can be dragged with the mouse.
type Scroll struct {
	offset        float32
	scroll        gesture.Scroll
	bar           kit.Scrollbar
	contentH      int
	viewportH     int
	content       func() elements.Element
	st            styles.Styles
	onReachBottom func()
	wasAtBottom   bool
	pxPerDp       float32
}

// New returns a Scroll around content, rebuilt fresh every frame. The
// styles argument is optional: Background is the container fill and
// ScrollBar customizes the scrollbar's width, colors, corner radius
// and minimum thumb length.
func New(content func() elements.Element, st ...styles.Styles) *Scroll {
	s := &Scroll{content: content}
	if len(st) > 0 {
		s.st = st[0]
	}
	return s
}

// ScrollTo sets the vertical offset in pixels.
func (s *Scroll) ScrollTo(v float32) {
	s.offset = v
}

// Offset returns the current vertical offset.
func (s *Scroll) Offset() float32 {
	return s.offset
}

// ContentHeight returns the total height of the scrolled content, as
// measured on the last frame.
func (s *Scroll) ContentHeight() int {
	return s.contentH
}

// ViewportHeight returns the height of the visible viewport, as
// measured on the last frame.
func (s *Scroll) ViewportHeight() int {
	return s.viewportH
}

// Progress returns how far scrolled, from 0 (top) to 1 (bottom); 0
// when there's nothing to scroll (the content fits the viewport).
func (s *Scroll) Progress() float32 {
	max := s.maxOffset()
	if max <= 0 {
		return 0
	}
	return s.offset / float32(max)
}

// AtBottom reports whether the scroll is at (or past) its maximum
// offset - also true when there's nothing to scroll at all (the
// content already fits the viewport), the same way a page with no
// scrollbar is trivially "at the bottom": there's nowhere else to go,
// so a load-more handler driven by AtBottom/OnReachBottom fires
// immediately, prompting an initial fetch when the first page of
// content doesn't even fill the viewport.
func (s *Scroll) AtBottom() bool {
	return s.offset >= float32(s.maxOffset())
}

// Visible reports whether the vertical range [top, top+height) - in
// the same coordinate space as Content's own output, e.g. a running
// offset tracked while building a long list - is at least partially
// within the current viewport. Useful for lazy-loading or skipping
// expensive work for an item that isn't currently on screen.
func (s *Scroll) Visible(top, height int) bool {
	return top+height > int(s.offset) && top < int(s.offset)+s.viewportH
}

// VisibleDp is Visible with top and height in dp, the units Styles
// use, converted with the display scale of the last frame.
func (s *Scroll) VisibleDp(top, height float32) bool {
	scale := s.pxPerDp
	if scale <= 0 {
		scale = 1
	}
	return s.Visible(int(top*scale), int(height*scale+0.5))
}

// OffsetDp returns the vertical offset in dp.
func (s *Scroll) OffsetDp() float32 {
	scale := s.pxPerDp
	if scale <= 0 {
		scale = 1
	}
	return s.offset / scale
}

// ScrollToDp sets the vertical offset in dp.
func (s *Scroll) ScrollToDp(v float32) {
	scale := s.pxPerDp
	if scale <= 0 {
		scale = elements.DisplayScale()
	}
	s.offset = v * scale
}

// OnReachBottom sets the handler invoked once when the scroll reaches
// its maximum offset - edge-triggered: it fires again only after
// scrolling away from the bottom and back, not on every frame spent
// there. It returns s for chaining.
func (s *Scroll) OnReachBottom(fn func()) *Scroll {
	s.onReachBottom = fn
	return s
}

// maxOffset returns the largest valid offset.
func (s *Scroll) maxOffset() int {
	m := s.contentH - s.viewportH
	if m < 0 {
		return 0
	}
	return m
}

// clamp bounds v to the valid offset range.
func (s *Scroll) clamp(v float32) float32 {
	maxOff := float32(s.maxOffset())
	if v < 0 {
		return 0
	}
	if v > maxOff {
		return maxOff
	}
	return v
}

// Layout lays out the scrollable content.
func (s *Scroll) Layout(gtx layout.Context) layout.Dimensions {
	viewport := gtx.Constraints.Constrain(image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y))
	gtx.Constraints = layout.Exact(viewport)
	s.viewportH = viewport.Y
	s.pxPerDp = gtx.Metric.PxPerDp

	// Lay the content out with unbounded height.
	macro := op.Record(gtx.Ops)
	cgtx := gtx
	cgtx.Constraints.Min.Y = 0
	cgtx.Constraints.Max.Y = unbounded
	track := elements.TrackChild()
	dims := s.content().Layout(cgtx)
	track.End()
	contentOps := macro.Stop()
	s.contentH = dims.Size.Y

	// Wheel and touch scrolling. The scroll region covers the
	// viewport and must be registered below the content so
	// interactive children keep priority.
	defer clip.Rect{Max: viewport}.Push(gtx.Ops).Pop()
	// Keep touch-drag scrolling working even inside a frameless
	// window's whole-window drag area - see slider.Layout's identical
	// comment.
	kit.BlockDragHere(gtx)
	s.scroll.Add(gtx.Ops)

	// Clamp the offset to the content.
	s.offset = s.clamp(s.offset)

	// The ranges accept any scroll distance that stays within the
	// content.
	scrollY := pointer.ScrollRange{
		Min: -int(s.offset),
		Max: s.maxOffset() - int(s.offset),
	}
	dist := s.scroll.Update(gtx.Metric, gtx.Source, gtx.Now, gesture.Vertical, pointer.ScrollRange{}, scrollY)
	if dist != 0 {
		s.offset = s.clamp(s.offset + float32(dist))
	}
	if dist != 0 || s.scroll.State() != gesture.StateIdle {
		// Flings keep scrolling between input events; cached
		// ancestors must keep laying the content out.
		kit.MarkLive()
	}

	atBottom := s.AtBottom()
	if atBottom && !s.wasAtBottom && s.onReachBottom != nil {
		s.onReachBottom()
	}
	s.wasAtBottom = atBottom

	// Draw the scrolled content.
	if properties.CalcColor(s.st.Background).A != 0 {
		paint.Fill(gtx.Ops, properties.CalcColor(s.st.Background))
	}
	// The content was laid out before the offset was final: move what
	// it registered (dropdown anchors, semantic nodes) with it.
	track.Place(image.Pt(0, -int(s.offset)))
	off := op.Offset(image.Pt(0, -int(s.offset))).Push(gtx.Ops)
	contentOps.Add(gtx.Ops)
	off.Pop()

	s.layoutScrollbar(gtx, viewport)
	return layout.Dimensions{Size: viewport}
}

// layoutScrollbar draws the track and draggable thumb when the
// content overflows, fully customized via st.ScrollBar.
func (s *Scroll) layoutScrollbar(gtx layout.Context, viewport image.Point) {
	s.bar.Layout(gtx, viewport, s.contentH, &s.offset, s.st.ScrollBar)
}
