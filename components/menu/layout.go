package menu

import (
	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/components/kit"
)

// Layout renders the open menu: a window-wide dismissal layer with
// the panel on top, clamped to the window bounds. Layout it last in
// the view so it renders above the rest of the UI.
func (m *Menu) Layout(gtx layout.Context) layout.Dimensions {
	if !m.open {
		return layout.Dimensions{}
	}

	// Any primary click outside the panel closes it.
	if m.scrim.Clicked(gtx) {
		m.open = false
		return layout.Dimensions{}
	}
	m.scrim.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: gtx.Constraints.Max}
	})

	// Escape closes the menu. Focus is requested below, before the
	// events of this frame are delivered.
	for {
		ev, ok := gtx.Event(key.Filter{Focus: m, Name: key.NameEscape})
		if !ok {
			break
		}
		if e, ok := ev.(key.Event); ok && e.State == key.Press {
			m.open = false
			return layout.Dimensions{}
		}
	}

	// Measure the panel, then position it within the window.
	macro := op.Record(gtx.Ops)
	dims := m.panel(gtx)
	panel := macro.Stop()

	win := f32.Pt(float32(gtx.Constraints.Max.X), float32(gtx.Constraints.Max.Y))
	pos := kit.ClampPos(win, f32.Pt(float32(dims.Size.X), float32(dims.Size.Y)), m.pos, kit.Margin)

	off := op.Offset(pos.Round()).Push(gtx.Ops)
	panel.Add(gtx.Ops)
	off.Pop()

	// Keep focus for the Escape handler while open.
	gtx.Execute(key.FocusCmd{Tag: m})
	return layout.Dimensions{Size: gtx.Constraints.Max}
}
