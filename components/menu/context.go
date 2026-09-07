package menu

import (
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
)

// contextTag identifies pointer events delivered to the global
// right-click area.
type contextTag struct {
	m *Menu
}

// Context registers a window-wide right-click area. It must be laid
// out with a context whose origin is the window origin (the app does
// this for menus passed to SetContextMenu) and before the rest of the
// UI so widgets keep priority over it. It never paints.
func (m *Menu) Context(gtx layout.Context) layout.Dimensions {
	tag := &contextTag{m: m}
	ps := pointer.PassOp{}.Push(gtx.Ops)
	cl := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	event.Op(gtx.Ops, tag)
	cl.Pop()
	ps.Pop()

	for {
		ev, ok := gtx.Source.Event(pointer.Filter{
			Target: tag,
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
