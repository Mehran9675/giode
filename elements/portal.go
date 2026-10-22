package elements

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"

	"github.com/mehran9675/giode/components/kit"
)

// portalEntry is one piece of content queued to render above
// everything else, regardless of how deeply nested its caller was.
type portalEntry struct {
	// anchor is the window position the content goes at, before it's
	// moved to stay inside the window. Containers that place children
	// after laying them out correct it (see ShiftPositions).
	anchor image.Point
	size   image.Point
	build  func() Element
	// dismiss, when set, is the tag of a full-window layer under the
	// content that closes the popup on a press outside it.
	dismiss *dismissLayer
	// fixed entries were given a final window position (RegisterPortalAt)
	// and aren't moved by containers.
	fixed bool
}

// dismissLayer is the outside-click layer of a popup. Its identity
// must be stable across frames, so it's kept per owner.
type dismissLayer struct {
	opts  PopupOptions
	shown bool
}

// PopupOptions configures RegisterPopup.
type PopupOptions struct {
	// OnDismiss runs when the user presses outside the popup.
	OnDismiss func()
	// PassThrough lets that press also reach whatever is under it,
	// instead of only closing the popup. Use it for popups attached to
	// a text field (suggestions), so clicking back into the field
	// still places the caret.
	PassThrough bool
}

var (
	portalQueue   []portalEntry
	dismissLayers = map[any]*dismissLayer{}
)

// windowSize is the window's size for the current frame, recorded by
// SetWindowSize. Portals are measured and kept inside it.
var windowSize image.Point

// SetWindowSize records the window's size for the current frame. The
// app calls this once, before laying out the view; component authors
// calling RegisterPortal don't need to call it themselves.
func SetWindowSize(size image.Point) {
	windowSize = size
}

// RegisterPortal queues build to be laid out after the rest of the
// view, on top of it, near anchor - the mechanism behind Select's
// list, menus, tooltips and AutoComplete's suggestions, and the way
// for any component to draw above everything else however deeply it's
// nested.
//
// anchor is a window position, normally AbsolutePosition() plus an
// offset such as the trigger's height:
//
//	anchor := elements.AbsolutePosition().Add(image.Pt(0, dims.Size.Y+4))
//	elements.RegisterPortal(gtx, anchor, func() elements.Element { return panel })
//
// It stays correct when the caller sits inside a Row, Stack, Scroll or
// other container that positions its children after laying them out:
// those containers move the anchor along with the child. When the
// content is drawn it's kept inside the window, flipping above the
// anchor when there's no room below.
func RegisterPortal(gtx layout.Context, anchor image.Point, build func() Element) {
	registerPortal(gtx, anchor, build, nil)
}

// RegisterPopup is RegisterPortal for content that should close when
// the user clicks anywhere else, like a dropdown list. A transparent
// layer covering the window is placed under the content; a press on it
// calls opts.OnDismiss. Unless opts.PassThrough is set, the press
// doesn't reach the UI below, the way native menus behave. owner
// identifies the popup across frames (usually the component itself)
// and must be comparable.
func RegisterPopup(gtx layout.Context, anchor image.Point, owner any, build func() Element, opts PopupOptions) {
	layer, ok := dismissLayers[owner]
	if !ok {
		layer = &dismissLayer{}
		dismissLayers[owner] = layer
	}
	layer.opts = opts
	registerPortal(gtx, anchor, build, layer)
}

// RegisterPortalAt is RegisterPortal for a position that is already a
// final window position rather than one computed from
// AbsolutePosition: the pointer position of a right-click, or an
// anchor kept with TrackPosition. Containers around the caller don't
// move it.
func RegisterPortalAt(gtx layout.Context, pos image.Point, build func() Element) {
	registerPortal(gtx, pos, build, nil)
	portalQueue[len(portalQueue)-1].fixed = true
}

func registerPortal(gtx layout.Context, anchor image.Point, build func() Element, dismiss *dismissLayer) {
	win := windowSize
	if win == (image.Point{}) {
		win = gtx.Constraints.Max
	}
	mgtx := gtx
	mgtx.Constraints = layout.Constraints{Max: win}
	dims := Measure(mgtx, build())
	portalQueue = append(portalQueue, portalEntry{anchor: anchor, size: dims.Size, build: build, dismiss: dismiss})
}

// place returns where an entry is drawn: its anchor, moved to keep the
// content inside the window.
func (e portalEntry) place(win image.Point) image.Point {
	winF := f32.Pt(float32(win.X), float32(win.Y))
	size := f32.Pt(float32(e.size.X), float32(e.size.Y))
	return kit.ClampPos(winF, size, f32.Pt(float32(e.anchor.X), float32(e.anchor.Y)), kit.Margin).Round()
}

// DrainPortals lays out and paints every portal registered so far this
// frame, in registration order (later ones on top), then clears the
// queue. The app calls this once per frame, after the rest of the view.
// Portals registered while draining (a submenu opened from its parent
// menu) are drained too.
func DrainPortals(gtx layout.Context) {
	win := windowSize
	if win == (image.Point{}) {
		win = gtx.Constraints.Max
	}
	for len(portalQueue) > 0 {
		entries := portalQueue
		portalQueue = nil
		for _, e := range entries {
			if e.dismiss != nil {
				layoutDismissLayer(gtx, e.dismiss, win)
			}
			pos := e.place(win)
			off := op.Offset(pos).Push(gtx.Ops)
			pop := PushOffset(pos)
			e.build().Layout(gtx)
			pop()
			off.Pop()
		}
	}
	// Forget the layers of popups that weren't shown this frame.
	for owner, layer := range dismissLayers {
		if !layer.shown {
			delete(dismissLayers, owner)
		}
		layer.shown = false
	}
}

// layoutDismissLayer covers the window with a press catcher that
// dismisses its popup.
func layoutDismissLayer(gtx layout.Context, l *dismissLayer, win image.Point) {
	l.shown = true
	for {
		ev, ok := gtx.Event(pointer.Filter{Target: l, Kinds: pointer.Press})
		if !ok {
			break
		}
		if pe, ok := ev.(pointer.Event); ok && pe.Kind == pointer.Press && l.opts.OnDismiss != nil {
			l.opts.OnDismiss()
			kit.RequestFrame(gtx)
		}
	}
	if l.opts.PassThrough {
		defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	}
	cl := clip.Rect{Max: win}.Push(gtx.Ops)
	event.Op(gtx.Ops, l)
	cl.Pop()
}
