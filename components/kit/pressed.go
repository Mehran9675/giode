package kit

import "image/color"

// Pressed returns the pressed shading of c: slightly darker.
func Pressed(c color.NRGBA) color.NRGBA {
	return Scale(c, 0.85)
}
