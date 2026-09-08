package input

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/fonts"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// editorEl renders the editor with a placeholder underneath.
type editorEl struct {
	input *Input
	st    styles.Styles
}

// Layout renders the editor.
func (e *editorEl) Layout(gtx layout.Context) layout.Dimensions {
	st := e.st
	size := properties.ResolveFontSize(st.FontSize)
	if e.input.editor.Text() == "" {
		gtx.Constraints.Min.X = 80
	}

	textMacro := op.Record(gtx.Ops)
	paint.ColorOp{Color: properties.ResolveColor(st.Color)}.Add(gtx.Ops)
	textColor := textMacro.Stop()

	selMacro := op.Record(gtx.Ops)
	sel := properties.ResolveColor(st.Color)
	sel.A = uint8(float32(sel.A) * 0.3)
	paint.ColorOp{Color: sel}.Add(gtx.Ops)
	selColor := selMacro.Stop()

	// The placeholder is painted first so the editor (caret and
	// selection) renders on top of it.
	if e.input.editor.Text() == "" && e.input.placeholder != "" {
		hintColor := properties.ResolveColor(st.Color)
		hintColor.A = uint8(float32(hintColor.A) * 0.5)
		elements.Text(e.input.placeholder, styles.Styles{Color: properties.CalcColorReverse(hintColor), FontSize: st.FontSize}).Layout(gtx)
	}

	return e.input.editor.Layout(gtx, fonts.Shaper(), font.Font{}, size, textColor, selColor)
}
