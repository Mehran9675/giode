package slider

import "image/color"

var (
	defaultTrack = color.NRGBA{R: 0x33, G: 0x41, B: 0x55, A: 0xff}
	defaultFill  = color.NRGBA{R: 0x3b, G: 0x82, B: 0xf6, A: 0xff}
	defaultThumb = color.NRGBA{R: 0xf8, G: 0xfa, B: 0xfc, A: 0xff}
)

const (
	height  = 24
	trackH  = 4
	thumbD  = 14
	padding = thumbD / 2
)

func resolve(c color.NRGBA, def color.NRGBA) color.NRGBA {
	if c.A == 0 {
		return def
	}
	return c
}
