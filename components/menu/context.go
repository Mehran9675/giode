package menu

import (
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
)

// Context registers a window-wide right-click area. It must be laid
// out with a context whose origin is the window origin (the app does
// this for menus passed to SetContextMenu) and after the rest of the
// UI so it isn't blocked by ordinary widgets. It never paints.
//
// Ordering matters here in a way that's easy to get backwards: this
// area uses pointer.PassOp so it never blocks widgets behind it, but
// PassOp does NOT make it immune to being blocked BY a non-pass-through
// widget in front of it (Gio's hit-test latches into "blocked" mode
// for the rest of a reverse scan once it meets one such widget, and
// doesn't unlatch just because a pass-through area follows). Almost
// every pixel of a real UI is normally covered by some ordinary
// (non-pass) widget — a button, a scrollable area's own drag region,
// and so on — so registering this area first (behind everything, as
// one might expect from "let widgets take priority") means it is
// blocked almost everywhere and the menu never opens. Registering it
// last (in front of everything) instead lets it see every press while
// still leaving normal clicks on the widgets behind it untouched,
// since pointer delivery isn't exclusive between overlapping matches.
//
// m itself is used as the event tag: Gio correlates a pointer event
// across frames by tag identity, so the tag must be a stable pointer
// reused every frame (like Menu's own key.FocusCmd{Tag: m} in
// layout.go) rather than a fresh value allocated per call.
func (m *Menu) Context(gtx layout.Context) layout.Dimensions {
	ps := pointer.PassOp{}.Push(gtx.Ops)
	cl := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	event.Op(gtx.Ops, m)
	cl.Pop()
	ps.Pop()

	for {
		ev, ok := gtx.Source.Event(pointer.Filter{
			Target: m,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}
		if e.Buttons.Contain(pointer.ButtonSecondary) {
			m.OpenAt(e.Position)
		}
	}
	return layout.Dimensions{}
}
