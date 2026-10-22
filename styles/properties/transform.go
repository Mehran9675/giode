package properties

import (
	"image"
	"math"

	"gioui.org/f32"
)

// Transform rotates, scales and translates an element when painting
// (CSS transform). Like in CSS it doesn't affect layout: siblings are
// placed as if the element were untransformed, while pointer input is
// transformed along with the painted output. The zero value is the
// identity.
type Transform struct {
	// Rotate is the clockwise rotation in degrees.
	Rotate float32
	// Scale scales both axes; ScaleX and ScaleY, when non-zero,
	// override it per axis. Zero means 1 (unscaled).
	Scale, ScaleX, ScaleY float32
	// TranslateX and TranslateY offset the element, in dp.
	TranslateX, TranslateY int
	// OriginX and OriginY place the rotation/scale origin as a
	// fraction of the element size. Zero for both means the center
	// (0.5, 0.5), like CSS; use a tiny value such as 0.0001 to pin the
	// top-left corner.
	OriginX, OriginY float32
}

// IsIdentity reports whether t leaves the element unchanged.
func (t Transform) IsIdentity() bool {
	return t.Rotate == 0 && t.scaleX() == 1 && t.scaleY() == 1 && t.TranslateX == 0 && t.TranslateY == 0
}

func (t Transform) scaleX() float32 {
	switch {
	case t.ScaleX != 0:
		return t.ScaleX
	case t.Scale != 0:
		return t.Scale
	}
	return 1
}

func (t Transform) scaleY() float32 {
	switch {
	case t.ScaleY != 0:
		return t.ScaleY
	case t.Scale != 0:
		return t.Scale
	}
	return 1
}

// Affine returns the transformation for an element of the given size.
// TranslateX/TranslateY must already be resolved to pixels.
func (t Transform) Affine(size image.Point) f32.Affine2D {
	ox, oy := t.OriginX, t.OriginY
	if ox == 0 && oy == 0 {
		ox, oy = 0.5, 0.5
	}
	origin := f32.Pt(float32(size.X)*ox, float32(size.Y)*oy)
	return f32.AffineId().
		Scale(origin, f32.Pt(t.scaleX(), t.scaleY())).
		Rotate(origin, t.Rotate*math.Pi/180).
		Offset(f32.Pt(float32(t.TranslateX), float32(t.TranslateY)))
}
