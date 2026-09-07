package kit

import "image/color"

// Hovered returns the hover highlight of c: slightly lighter.
func Hovered(c color.NRGBA) color.NRGBA {
	return Scale(c, 1.15)
}
