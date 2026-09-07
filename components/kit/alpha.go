package kit

import "image/color"

// Alpha returns c with its alpha multiplied by a.
func Alpha(c color.NRGBA, a float32) color.NRGBA {
	return color.NRGBA{R: c.R, G: c.G, B: c.B, A: uint8(float32(c.A) * a)}
}
