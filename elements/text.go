package elements

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/fonts"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// TextElement renders a single string with its styles. Chain OnHover
// to add hover styles.
type TextElement struct {
	str string
	st  styles.Styles

	kit.HoverTracker
	hoverSt styles.Styles
}

// Text renders a string with the given styles. The styles argument is
// optional.
func Text(str string, st ...styles.Styles) *TextElement {
	var s styles.Styles
	if len(st) > 0 {
		s = st[0]
	}
	return &TextElement{str: str, st: s}
}

// OnHover sets styles applied while the pointer is over the text,
// merged on top of its base styles - color, decoration, weight,
// anything Styles' text properties support:
//
//	label := giode.Text("Docs", giode.Styles{Color: "#94a3b8"}).OnHover(giode.Styles{
//		Color:          "#f8fafc",
//		TextDecoration: giode.TextDecorationUnderline,
//		Cursor:         giode.CursorPointer,
//	})
//
// Hover state lives in the element, so create it once and reuse it
// every frame; rebuilding it in the view function resets the hover on
// each frame.
func (t *TextElement) OnHover(st styles.Styles) *TextElement {
	t.Enable()
	t.hoverSt = st
	return t
}

func (t *TextElement) flexWeight(axis layout.Axis) int {
	return weightOf(t.st, axis)
}

func (t *TextElement) percentOf(axis layout.Axis) int {
	return percentOf(t.st, axis)
}

func (t *TextElement) zIndex() int {
	return int(t.st.ZIndex)
}

// resolveStyles returns the styles to lay out and paint with: the
// base styles, with the hover styles merged on top while hovered.
func (t *TextElement) resolveStyles() styles.Styles {
	if !t.Enabled() || !t.Hovered() {
		return t.st
	}
	return styles.Merge(t.st, t.hoverSt)
}

// Layout shapes and paints the text with its decoration.
func (t *TextElement) Layout(gtx layout.Context) layout.Dimensions {
	t.Update(gtx)
	st := t.resolveStyles()

	if fonts.Shaper() == nil {
		return layout.Dimensions{}
	}
	col := properties.ResolveColor(st.Color)
	colMacro := op.Record(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	textColor := colMacro.Stop()

	face := font.Font{
		Weight:   st.FontWeight.Weight(),
		Style:    st.FontStyle.Style(),
		Typeface: font.Typeface(st.FontFamily),
	}
	dims := widget.Label{
		Alignment:  st.TextAlign.LayoutAlignment(),
		MaxLines:   int(st.MaxLines),
		LineHeight: unit.Sp(st.LineHeight),
	}.Layout(gtx, fonts.Shaper(), face, properties.ResolveFontSize(st.FontSize), t.str, textColor)

	paintDecoration(gtx, dims.Size, st.TextDecoration, col)
	t.RegisterArea(gtx, image.Pt(0, 0), dims.Size, st.Cursor)
	return dims
}

// paintDecoration draws underline or line-through decorations.
func paintDecoration(gtx layout.Context, size image.Point, deco properties.TextDecoration, col color.NRGBA) {
	if deco == properties.TextDecorationNone || deco == "" {
		return
	}
	var line image.Rectangle
	switch deco {
	case properties.TextDecorationLineThrough:
		line = image.Rect(0, size.Y/2, size.X, size.Y/2+1)
	default:
		line = image.Rect(0, size.Y-1, size.X, size.Y)
	}
	if line.Dx() <= 0 {
		return
	}
	defer clip.Rect(line).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
