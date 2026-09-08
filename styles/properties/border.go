package properties

import (
	"image"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// BorderWidth is the width of the border drawn around an element, in
// pixels. Zero means no border.
type BorderWidth int

// BorderColor is the color of the border. A zero alpha color means no
// border is drawn.
type BorderColor = string

// PaintBorder strokes a border around an element of the given size.
func PaintBorder(ops *op.Ops, size image.Point, w BorderWidth, c BorderColor, radius BorderRadius) {
	if w <= 0 || CalcColor(c).A == 0 || size.X <= 0 || size.Y <= 0 {
		return
	}
	half := int(w) / 2
	rect := image.Rect(half, half, size.X-half, size.Y-half)
	if rect.Dx() <= 0 || rect.Dy() <= 0 {
		return
	}
	rr := clip.UniformRRect(rect, int(radius))
	st := clip.Stroke{Path: rr.Path(ops), Width: float32(w)}.Op().Push(ops)
	paint.ColorOp{Color: CalcColor(c)}.Add(ops)
	paint.PaintOp{}.Add(ops)
	st.Pop()
}
