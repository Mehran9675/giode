package properties

import (
	"image"
	"math"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// ColorStop is one color of a gradient at Offset, in the range 0-1.
type ColorStop struct {
	Offset float32
	Color  Color
}

// LinearGradient is a CSS linear-gradient(): colors blend along a line
// at Angle through the element's center. Set it as
// Styles.BackgroundGradient; it paints above Background and below
// BackgroundImage.
type LinearGradient struct {
	// Angle is the CSS gradient angle in degrees: 90 runs left to
	// right, 180 top to bottom, 270 right to left and 360 bottom to
	// top. Zero means the CSS default, 180.
	Angle float32
	// Stops are the colors along the line, in ascending Offset order.
	// At least two are needed.
	Stops []ColorStop
}

// Linear returns a gradient at angle degrees with colors spread evenly
// across it:
//
//	properties.Linear(90, "#3b82f6", "#8b5cf6") // left to right
func Linear(angle float32, colors ...Color) *LinearGradient {
	g := &LinearGradient{Angle: angle}
	for i, c := range colors {
		off := float32(0)
		if len(colors) > 1 {
			off = float32(i) / float32(len(colors)-1)
		}
		g.Stops = append(g.Stops, ColorStop{Offset: off, Color: c})
	}
	return g
}

// PaintGradient fills the current clip, a box of the given size, with
// g. Any number of stops is supported: Gio's gradient op blends two
// colors, so every further segment is painted over the half-plane past
// its own start.
func PaintGradient(ops *op.Ops, size image.Point, g *LinearGradient) {
	if g == nil || len(g.Stops) < 2 || size.X <= 0 || size.Y <= 0 {
		return
	}
	angle := float64(g.Angle)
	if angle == 0 {
		angle = 180
	}
	rad := angle * math.Pi / 180
	dir := f32.Pt(float32(math.Sin(rad)), float32(-math.Cos(rad)))
	w, h := float64(size.X), float64(size.Y)
	length := float32(math.Abs(w*math.Sin(rad)) + math.Abs(h*math.Cos(rad)))
	center := f32.Pt(float32(w)/2, float32(h)/2)
	start := center.Sub(dir.Mul(length / 2))
	at := func(off float32) f32.Point { return start.Add(dir.Mul(length * off)) }
	perp := f32.Pt(-dir.Y, dir.X)
	big := float32(w+h) * 2

	for i := 0; i < len(g.Stops)-1; i++ {
		s1, s2 := g.Stops[i], g.Stops[i+1]
		p1, p2 := at(s1.Offset), at(s2.Offset)
		if p1 == p2 {
			p2 = p2.Add(dir.Mul(0.01))
		}
		var stack clip.Stack
		if i > 0 {
			var p clip.Path
			p.Begin(ops)
			p.MoveTo(p1.Add(perp.Mul(big)))
			p.LineTo(p1.Sub(perp.Mul(big)))
			p.LineTo(p1.Sub(perp.Mul(big)).Add(dir.Mul(big)))
			p.LineTo(p1.Add(perp.Mul(big)).Add(dir.Mul(big)))
			p.Close()
			stack = clip.Outline{Path: p.End()}.Op().Push(ops)
		}
		paint.LinearGradientOp{
			Stop1: p1, Color1: CalcColor(s1.Color),
			Stop2: p2, Color2: CalcColor(s2.Color),
		}.Add(ops)
		paint.PaintOp{}.Add(ops)
		if i > 0 {
			stack.Pop()
		}
	}
}
