package properties

import (
	"image"

	"gioui.org/op"
	"gioui.org/op/paint"
)

// BorderWidth is the width of the border drawn around an element, in
// dp. Zero means no border.
type BorderWidth int

// BorderColor is the color of the border. A zero alpha color means no
// border is drawn.
type BorderColor = string

// PaintBorder strokes a border around an element of the given size.
func PaintBorder(ops *op.Ops, size image.Point, w BorderWidth, c BorderColor, radius BorderRadius) {
	PaintBorderCorners(ops, size, w, c, UniformCorners(int(radius)))
}

// PaintBorderCorners strokes a border around an element of the given
// size, rounding each corner by its own radius. The stroke sits inside
// the box, so a border never paints outside the element's bounds.
func PaintBorderCorners(ops *op.Ops, size image.Point, w BorderWidth, c BorderColor, corners Corners) {
	col := CalcColor(c)
	if w <= 0 || col.A == 0 || size.X <= 0 || size.Y <= 0 {
		return
	}
	StrokeRRect(ops, image.Rectangle{Max: size}, corners, int(w), false, func() {
		paint.ColorOp{Color: col}.Add(ops)
		paint.PaintOp{}.Add(ops)
	})
}
