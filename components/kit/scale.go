// Package kit provides the small building blocks used to create
// custom components: color state shading and overlay positioning.
//
// See the "Extending the library" documentation for a full recipe.
package kit

import "image/color"

// Scale multiplies each channel of c by f, clamping to the valid
// range.
func Scale(c color.NRGBA, f float32) color.NRGBA {
	sc := func(v uint8) uint8 {
		n := int(float32(v) * f)
		if n > 0xff {
			n = 0xff
		}
		return uint8(n)
	}
	return color.NRGBA{R: sc(c.R), G: sc(c.G), B: sc(c.B), A: c.A}
}
