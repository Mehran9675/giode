// Package scroll provides a stateful vertical scroll container.
package scroll

import (
	"image"
	"image/color"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
)

var (
	defaultTrack = color.NRGBA{R: 0x47, G: 0x55, B: 0x69, A: 0x99}
	defaultThumb = color.NRGBA{R: 0x94, G: 0xa3, B: 0xb8, A: 0xcc}
)

const (
	defaultBarWidth = 8
	defaultThumbLen = 20

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
	offset    float32
	scroll    gesture.Scroll
	barDrag   gesture.Drag
	barDragY  float32
	contentH  int
	viewportH int
	content   func() elements.Element
	st        styles.Styles
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

	// Lay the content out with unbounded height.
	macro := op.Record(gtx.Ops)
	cgtx := gtx
	cgtx.Constraints.Min.Y = 0
	cgtx.Constraints.Max.Y = unbounded
	dims := s.content().Layout(cgtx)
	contentOps := macro.Stop()
	s.contentH = dims.Size.Y

	// Wheel and touch scrolling. The scroll region covers the
	// viewport and must be registered below the content so
	// interactive children keep priority.
	defer clip.Rect{Max: viewport}.Push(gtx.Ops).Pop()
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

	// Draw the scrolled content.
	if properties.CalcColor(s.st.Background).A != 0 {
		paint.Fill(gtx.Ops, properties.CalcColor(s.st.Background))
	}
	off := op.Offset(image.Pt(0, -int(s.offset))).Push(gtx.Ops)
	contentOps.Add(gtx.Ops)
	off.Pop()

	s.layoutScrollbar(gtx, viewport)
	return layout.Dimensions{Size: viewport}
}

// layoutScrollbar draws the track and draggable thumb when the
// content overflows, fully customized via st.ScrollBar.
func (s *Scroll) layoutScrollbar(gtx layout.Context, viewport image.Point) {
	if s.contentH <= s.viewportH || s.viewportH <= 0 {
		return
	}
	sb := s.st.ScrollBar
	barWidth := sb.Width
	if barWidth <= 0 {
		barWidth = defaultBarWidth
	}
	minThumb := sb.MinThumbLength
	if minThumb <= 0 {
		minThumb = defaultThumbLen
	}
	radius := sb.Radius
	if radius <= 0 {
		radius = barWidth / 2
	}
	trackColor := properties.CalcColor(sb.TrackColor)
	if trackColor.A == 0 {
		trackColor = defaultTrack
	}
	thumbColor := properties.CalcColor(sb.ThumbColor)
	if thumbColor.A == 0 {
		thumbColor = defaultThumb
	}

	x := viewport.X - barWidth
	trackH := viewport.Y
	thumbH := trackH * trackH / s.contentH
	if thumbH < minThumb {
		thumbH = minThumb
	}
	maxThumbOff := trackH - thumbH
	thumbOff := int(s.offset) * maxThumbOff / max(s.contentH-s.viewportH, 1)

	// Thumb drag region and events.
	thumbRect := image.Rect(x, thumbOff, x+barWidth, min(thumbOff+thumbH, trackH))
	thumbClip := clip.Rect(thumbRect).Push(gtx.Ops)
	s.barDrag.Add(gtx.Ops)
	for {
		ev, ok := s.barDrag.Update(gtx.Metric, gtx.Source, gesture.Both)
		if !ok {
			break
		}
		switch ev.Kind {
		case pointer.Press:
			if ev.Priority == pointer.Grabbed {
				s.barDragY = ev.Position.Y
			}
		case pointer.Drag:
			if ev.Priority == pointer.Grabbed {
				dy := ev.Position.Y - s.barDragY
				s.barDragY = ev.Position.Y
				if dy != 0 && maxThumbOff > 0 {
					delta := dy * float32(s.contentH-s.viewportH) / float32(maxThumbOff)
					s.offset = s.clamp(s.offset + delta)
				}
			}
		}
	}
	thumbClip.Pop()

	// Track and thumb, positioned at the same x offset as the drag
	// region above.
	barOff := op.Offset(image.Pt(x, 0)).Push(gtx.Ops)
	trackClip := clip.UniformRRect(image.Rect(0, 0, barWidth, trackH), radius).Push(gtx.Ops)
	paint.Fill(gtx.Ops, trackColor)
	trackClip.Pop()

	thumbOff = int(s.offset) * maxThumbOff / max(s.contentH-s.viewportH, 1)
	thumbClip2 := clip.UniformRRect(image.Rect(0, thumbOff, barWidth, thumbOff+thumbH), radius).Push(gtx.Ops)
	paint.Fill(gtx.Ops, thumbColor)
	thumbClip2.Pop()
	barOff.Pop()
}
