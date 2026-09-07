package dialog

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles/properties"
)

// Layout renders the dialog: a dimmed full-window scrim with the
// centered panel. A click on the scrim closes it.
func (d *Dialog) Layout(gtx layout.Context, content elements.Element) layout.Dimensions {
	if !d.open {
		return layout.Dimensions{}
	}
	win := gtx.Constraints.Max

	if d.scrim.Clicked(gtx) {
		d.open = false
		return layout.Dimensions{}
	}
	d.scrim.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		defer clip.Rect{Max: win}.Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, color.NRGBA{A: 0x66})
		return layout.Dimensions{Size: win}
	})

	st := d.st
	if st.Background.A == 0 {
		st.Background = defaultBackground
	}
	if st.BorderRadius == 0 {
		st.BorderRadius = 12
	}
	if st.Padding == (properties.Inset{}) {
		st.Padding = properties.UniformInset(20)
	}

	// Measure the panel, then center it.
	macro := op.Record(gtx.Ops)
	panelGtx := gtx
	panelGtx.Constraints.Max.X = win.X - 2*kit.Margin
	panelGtx.Constraints.Max.Y = win.Y - 2*kit.Margin
	dims := elements.Box(st, content).Layout(panelGtx)
	panel := macro.Stop()

	x := (win.X - dims.Size.X) / 2
	y := (win.Y - dims.Size.Y) / 2
	off := op.Offset(image.Pt(x, y)).Push(gtx.Ops)
	panel.Add(gtx.Ops)
	off.Pop()

	return layout.Dimensions{Size: win}
}
