package properties

import (
	"image"
	"image/color"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// Background is the fill painted behind an element's content (CSS
// background-color). A zero alpha background is transparent and not
// painted.
type Background = color.NRGBA

// PaintBackground fills size with bg, optionally with rounded corners.
func PaintBackground(ops *op.Ops, size image.Point, bg Background, radius BorderRadius) {
	if bg.A == 0 || size.X <= 0 || size.Y <= 0 {
		return
	}
	defer clip.UniformRRect(image.Rectangle{Max: size}, int(radius)).Push(ops).Pop()
	paint.Fill(ops, bg)
}
