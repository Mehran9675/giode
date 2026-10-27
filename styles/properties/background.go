package properties

import (
	"image"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// Background is the fill painted behind an element's content (CSS
// background-color). A zero alpha background is transparent and not
// painted.
type Background = string

// BackgroundGradient is a gradient painted behind an element's content,
// above Background (CSS background-image: linear-gradient(...)).
type BackgroundGradient = *LinearGradient

// PaintBackground fills size with bg, optionally with rounded corners.
func PaintBackground(ops *op.Ops, size image.Point, bg Background, radius BorderRadius) {
	PaintBackgroundCorners(ops, size, bg, UniformCorners(int(radius)))
}

// PaintBackgroundCorners fills size with bg, rounding each corner by
// its own radius.
func PaintBackgroundCorners(ops *op.Ops, size image.Point, bg Background, c Corners) {
	col := CalcColor(bg)
	if col.A == 0 || size.X <= 0 || size.Y <= 0 {
		return
	}
	FillRRect(ops, image.Rectangle{Max: size}, c, false, func() {
		paint.ColorOp{Color: col}.Add(ops)
		paint.PaintOp{}.Add(ops)
	})
}

// PaintGradientCorners fills size with g, rounding each corner by its
// own radius.
func PaintGradientCorners(ops *op.Ops, size image.Point, g *LinearGradient, c Corners) {
	if g == nil || size.X <= 0 || size.Y <= 0 {
		return
	}
	FillRRect(ops, image.Rectangle{Max: size}, c, false, func() { PaintGradient(ops, size, g) })
}

// RRect returns the rounded rectangle clip for r with per-corner
// radii, clamped so they never exceed half the rectangle's size.
func RRect(r image.Rectangle, c Corners) clip.RRect {
	c = clampCorners(c, r.Dx(), r.Dy())
	return clip.RRect{Rect: r, NW: c.TopLeft, NE: c.TopRight, SE: c.BottomRight, SW: c.BottomLeft}
}
