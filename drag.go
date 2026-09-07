package giode

import (
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op/clip"
)

// dragState tracks whether the pointer is over empty space (no widget
// claimed it) and, while so, registers the whole window as a native
// move area so the OS drags the window from its background.
type dragState struct {
	active bool
}

// layout registers the pointer region and, when the pointer hovers
// empty space, the native move area.
func (d *dragState) layout(gtx layout.Context) layout.Dimensions {
	// A pass-through region sees every pointer event but never blocks
	// the widgets above: Grabbed priority means the pointer is over
	// empty space, Shared means a widget claimed it.
	ps := pointer.PassOp{}.Push(gtx.Ops)
	cl := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	event.Op(gtx.Ops, d)
	cl.Pop()
	ps.Pop()

	for {
		ev, ok := gtx.Source.Event(pointer.Filter{
			Target: d,
			Kinds:  pointer.Move | pointer.Press | pointer.Release | pointer.Drag | pointer.Cancel,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}
		d.active = e.Priority == pointer.Grabbed
	}

	if d.active {
		// The action area makes the OS treat presses as caption
		// drags. It is registered only while the pointer hovers
		// empty space, so widgets stay clickable.
		system.ActionInputOp(system.ActionMove).Add(gtx.Ops)
	}
	return layout.Dimensions{}
}
