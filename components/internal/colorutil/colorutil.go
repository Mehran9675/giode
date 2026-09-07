// Package colorutil provides small color helpers shared by the
// components.
package colorutil

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

// Hovered returns the hover highlight of c.
func Hovered(c color.NRGBA) color.NRGBA {
	return Scale(c, 1.15)
}

// Pressed returns the pressed shading of c.
func Pressed(c color.NRGBA) color.NRGBA {
	return Scale(c, 0.85)
}

// Alpha returns c with its alpha multiplied by a.
func Alpha(c color.NRGBA, a float32) color.NRGBA {
	return color.NRGBA{R: c.R, G: c.G, B: c.B, A: uint8(float32(c.A) * a)}
}
