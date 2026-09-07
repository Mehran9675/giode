package properties

import (
	"image"
	"image/color"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// BoxShadow is a drop shadow painted behind an element (CSS
// box-shadow). A zero alpha color means no shadow.
type BoxShadow struct {
	// X and Y offset the shadow from the element.
	X, Y int
	// Blur is the shadow radius in pixels. Zero gives a crisp edge.
	Blur int
	// Color is the shadow color; its alpha scales the shadow.
	Color color.NRGBA
}

// PaintBoxShadow paints the shadow for an element whose painted box
// is rect with the given corner radius.
func PaintBoxShadow(ops *op.Ops, rect image.Rectangle, radius BorderRadius, s BoxShadow) {
	if s.Color.A == 0 || rect.Dx() <= 0 || rect.Dy() <= 0 {
		return
	}
	r := rect.Add(image.Pt(s.X, s.Y))
	blur := s.Blur
	if blur <= 0 {
		strokeRoundedRect(ops, r.Inset(-1), int(radius), s.Color)
		return
	}
	if blur > 16 {
		blur = 16
	}
	// Layer strokes from the outside in with a quadratic alpha falloff
	// for a soft edge.
	for i := blur; i > 0; i-- {
		f := float32(i) / float32(blur+1)
		c := s.Color
		c.A = uint8(float32(c.A) * (1 - f) * (1 - f))
		strokeRoundedRect(ops, r.Inset(-i/2), int(radius)+i, c)
	}
}

func strokeRoundedRect(ops *op.Ops, rect image.Rectangle, radius int, c color.NRGBA) {
	rr := clip.UniformRRect(rect, radius)
	st := clip.Stroke{Path: rr.Path(ops), Width: 1}.Op().Push(ops)
	paint.ColorOp{Color: c}.Add(ops)
	paint.PaintOp{}.Add(ops)
	st.Pop()
}
