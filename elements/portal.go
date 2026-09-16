package elements

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/components/kit"
)

// portalEntry is one piece of content queued to render above
// everything else, regardless of how deeply nested its caller was.
type portalEntry struct {
	pos   image.Point
	build func() Element
}

var portalQueue []portalEntry

// windowSize is the window's root constraints for the current frame,
// recorded by SetWindowSize before the view is laid out.
// RegisterPortal measures and positions against it instead of its
// caller's own ambient constraints, which shrink the deeper a portal
// trigger sits in the tree (Tooltip, Dropdown and Menu's nested
// submenus are typically nested arbitrarily deep, unlike Menu's own
// right-click panel, which happens to be laid out directly off the
// window) - using the wrong, smaller bound wouldn't crash (the portal
// still paints in the right place, at DrainPortals' own root
// position), but would measure/clamp against a box smaller than the
// real window, positioning it more conservatively than necessary.
var windowSize image.Point

// SetWindowSize records the window's size for the current frame. The
// app calls this once, before laying out the view; component authors
// calling RegisterPortal don't need to call it themselves.
func SetWindowSize(size image.Point) {
	windowSize = size
}

// RegisterPortal queues build to be laid out and painted for real
// after the rest of the view - the mechanism behind Tooltip, nested
// Menu submenus and AutoComplete's suggestion panel, and the building
// block for any component that needs to render on top of everything
// else no matter how deeply it's nested in the tree.
//
// anchor is a window-absolute point (typically AbsolutePosition(),
// optionally offset - e.g. by the trigger element's own height) that
// the content is positioned near; it's first laid out once on a
// disabled scratch context, bounded by the window (see
// SetWindowSize) rather than gtx's own ambient constraints, to
// measure its size (the same technique Menu's own panel-width
// measurement uses) so the position can be clamped to stay inside the
// window, flipping above the anchor if it would overflow the bottom
// edge - see kit.ClampPos. build is then laid out for real, exactly
// once, when the queue is drained.
func RegisterPortal(gtx layout.Context, anchor image.Point, build func() Element) {
	win := windowSize
	if win == (image.Point{}) {
		win = gtx.Constraints.Max
	}

	scratch := gtx
	scratch.Constraints = layout.Constraints{Max: win}
	var sOps op.Ops
	scratch.Ops = &sOps
	scratch.Source = gtx.Source.Disabled()
	dims := build().Layout(scratch)

	winF := f32.Pt(float32(win.X), float32(win.Y))
	size := f32.Pt(float32(dims.Size.X), float32(dims.Size.Y))
	pos := kit.ClampPos(winF, size, f32.Pt(float32(anchor.X), float32(anchor.Y)), kit.Margin)

	portalQueue = append(portalQueue, portalEntry{pos: pos.Round(), build: build})
}

// DrainPortals lays out and paints every portal registered so far
// this frame, in registration order (later registrations paint on
// top of earlier ones), then clears the queue. The app calls this
// once per frame, after the rest of the view.
//
// AbsolutePosition/PushOffset are kept in sync with the paint offset
// here (not just op.Offset) so that content built inside a portal -
// nested Menu submenus register a further portal of their own,
// anchored off their own row's position - reports its true screen
// position rather than whatever was accumulated before the portal
// that contains it was originally triggered.
func DrainPortals(gtx layout.Context) {
	entries := portalQueue
	portalQueue = nil
	for _, e := range entries {
		off := op.Offset(e.pos).Push(gtx.Ops)
		pop := PushOffset(e.pos)
		e.build().Layout(gtx)
		pop()
		off.Pop()
	}
}
