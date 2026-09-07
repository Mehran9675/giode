package spinner

import (
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// paintFill fills the current clip with c.
func paintFill(ops *op.Ops, c color.NRGBA) {
	paint.ColorOp{Color: c}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

// clipArc builds a stroked circular arc using the same recipe as
// Gio's material loader.
func clipArc(ops *op.Ops, startAngle, endAngle, radius float32) clip.Op {
	const thickness = .25

	var (
		width = radius * thickness
		delta = endAngle - startAngle

		vy, vx = math.Sincos(float64(startAngle))

		inner  = radius * (1. - thickness*.5)
		pen    = f32.Pt(float32(vx), float32(vy)).Mul(inner)
		center = f32.Pt(0, 0).Sub(pen)

		p clip.Path
	)

	p.Begin(ops)
	p.Move(pen)
	p.Arc(center, center, delta)
	return clip.Stroke{
		Path:  p.End(),
		Width: width,
	}.Op()
}
