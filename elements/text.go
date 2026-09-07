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

	"github.com/mehran9675/giode/fonts"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// textEl renders a single string with its styles.
type textEl struct {
	str string
	st  styles.Styles
}

// Text renders a string with the given styles. The styles argument is
// optional.
func Text(str string, st ...styles.Styles) Element {
	var s styles.Styles
	if len(st) > 0 {
		s = st[0]
	}
	return &textEl{str: str, st: s}
}

func (t *textEl) flexWeight(axis layout.Axis) int {
	return weightOf(t.st, axis)
}

func (t *textEl) percentOf(axis layout.Axis) int {
	return percentOf(t.st, axis)
}

func (t *textEl) zIndex() int {
	return int(t.st.ZIndex)
}

// Layout shapes and paints the text with its decoration.
func (t *textEl) Layout(gtx layout.Context) layout.Dimensions {
	if t.st.Display == properties.None || fonts.Shaper() == nil {
		return layout.Dimensions{}
	}
	col := properties.ResolveColor(t.st.Color)
	colMacro := op.Record(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	textColor := colMacro.Stop()

	face := font.Font{
		Weight:   t.st.FontWeight.Weight(),
		Style:    t.st.FontStyle.Style(),
		Typeface: font.Typeface(t.st.FontFamily),
	}
	dims := widget.Label{
		Alignment:  t.st.TextAlign.LayoutAlignment(),
		MaxLines:   int(t.st.MaxLines),
		LineHeight: unit.Sp(t.st.LineHeight),
	}.Layout(gtx, fonts.Shaper(), face, properties.ResolveFontSize(t.st.FontSize), t.str, textColor)

	paintDecoration(gtx, dims.Size, t.st.TextDecoration, col)
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
