package menu

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/elements"
)

// Trigger wraps el so clicking it toggles the menu, anchored just
// below it - a widget-anchored, dropdown-style trigger, alongside
// SetContextMenu (right-click anywhere) and OpenAt/Toggle (any other
// custom trigger, e.g. a long-press).
//
// The returned element must be laid out somewhere in the view like el
// would be on its own; m itself must also still be laid out somewhere
// in the tree (its panel renders via elements.RegisterPortal
// regardless of where) for its open state to actually render.
func (m *Menu) Trigger(el elements.Element) elements.Element {
	return &menuTrigger{menu: m, trigger: el, click: new(widget.Clickable)}
}

type menuTrigger struct {
	menu    *Menu
	trigger elements.Element
	click   *widget.Clickable
	anchor  image.Point
}

func (t *menuTrigger) Layout(gtx layout.Context) layout.Dimensions {
	if t.click.Clicked(gtx) {
		if t.menu.open {
			t.menu.Close()
		} else {
			t.menu.open = true
		}
	}
	dims := t.click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return t.trigger.Layout(gtx)
	})
	if t.menu.open {
		// The menu opens below the trigger, at the position the
		// trigger had once its enclosing containers placed it on the
		// last layout.
		t.menu.pos = f32.Pt(float32(t.anchor.X), float32(t.anchor.Y))
	}
	elements.TrackPosition(&t.anchor, elements.AbsolutePosition().Add(image.Pt(0, dims.Size.Y+4)))
	// Keep the trigger clickable even when an ancestor is a frameless
	// window's whole-window drag area - see button.Layout's identical
	// comment.
	kit.BlockDrag(gtx, dims.Size)
	return dims
}
