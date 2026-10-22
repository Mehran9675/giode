package properties

import (
	"image"
	"image/color"
	"math"

	"gioui.org/op"
	"gioui.org/op/paint"
)

// BoxShadow is a drop shadow painted behind an element (CSS
// box-shadow). A zero alpha color means no shadow. Lengths are in dp.
type BoxShadow struct {
	// X and Y offset the shadow from the element.
	X, Y int
	// Blur is the shadow's soft-edge radius. Zero gives a crisp edge.
	Blur int
	// Spread grows (or, when negative, shrinks) the shadow before it
	// is blurred.
	Spread int
	// Color is the shadow color; its alpha scales the shadow.
	Color string
}

// maxShadowLayers bounds the number of layers a blurred shadow is
// built from, keeping large blurs cheap.
const maxShadowLayers = 24

// PaintBoxShadow paints the shadow for an element whose painted box
// is rect with the given corner radius.
func PaintBoxShadow(ops *op.Ops, rect image.Rectangle, radius BorderRadius, s BoxShadow) {
	PaintBoxShadowCorners(ops, rect, UniformCorners(int(radius)), s)
}

// PaintBoxShadowCorners paints the shadow for an element whose painted
// box is rect with per-corner radii. The blur is approximated with
// stacked translucent rounded rectangles whose alphas compose to the
// shadow color at the core and fade out linearly over Blur.
func PaintBoxShadowCorners(ops *op.Ops, rect image.Rectangle, corners Corners, s BoxShadow) {
	col := CalcColor(s.Color)
	if col.A == 0 || rect.Dx() <= 0 || rect.Dy() <= 0 {
		return
	}
	r := rect.Add(image.Pt(s.X, s.Y)).Inset(-s.Spread)
	if r.Dx() <= 0 || r.Dy() <= 0 {
		return
	}
	grow := func(c Corners, d int) Corners {
		g := func(v int) int {
			if v == 0 {
				return 0
			}
			return max(0, v+d)
		}
		return Corners{g(c.TopLeft), g(c.TopRight), g(c.BottomRight), g(c.BottomLeft)}
	}
	if s.Blur <= 0 {
		fill(ops, r, grow(corners, s.Spread), col)
		return
	}
	layers := min(s.Blur, maxShadowLayers)
	step := float32(s.Blur) / float32(layers)
	// Each layer's alpha is chosen so all layers together compose to
	// the requested alpha where they overlap.
	a := float64(col.A) / 255
	per := 1 - math.Pow(1-a, 1/float64(layers))
	lc := col
	lc.A = uint8(per*255 + 0.5)
	if lc.A == 0 {
		lc.A = 1
	}
	for i := 0; i < layers; i++ {
		// Layers span from -Blur/2 (inside) to +Blur/2 (outside) of
		// the shadow edge, like a Gaussian centered on it.
		d := int(float32(i)*step - float32(s.Blur)/2 + 0.5)
		lr := r.Inset(-d)
		if lr.Dx() <= 0 || lr.Dy() <= 0 {
			continue
		}
		fill(ops, lr, grow(corners, s.Spread+d), lc)
	}
}

func fill(ops *op.Ops, r image.Rectangle, c Corners, col color.NRGBA) {
	st := RRect(r, c).Push(ops)
	paint.ColorOp{Color: col}.Add(ops)
	paint.PaintOp{}.Add(ops)
	st.Pop()
}
