package textrender

import (
	"image/color"
	"math"
	"sync"
)

// Glyph coverage is corrected for the GPU's blending, which happens in
// linear light: blended there, coverage makes dark text thin and gray
// and light text heavy and soft, unlike the platforms' text, which is
// blended closer to the display's gamma. The correction raises
// coverage toward full for dark text and lowers it for light text, by
// the exponents below, mixed by the text color's luminance. Light text
// gets the milder one: more thins small light text until it breaks up.
const (
	darkTextGamma  = 1.8
	lightTextGamma = 1.3
)

var luts sync.Map // color.NRGBA (opaque) -> *[256]uint8

// coverageLUT maps a glyph's coverage to the alpha it's drawn with in
// col.
func coverageLUT(col color.NRGBA) *[256]uint8 {
	key := color.NRGBA{R: col.R, G: col.G, B: col.B, A: 0xff}
	if l, ok := luts.Load(key); ok {
		return l.(*[256]uint8)
	}
	lum := (0.2126*float64(col.R) + 0.7152*float64(col.G) + 0.0722*float64(col.B)) / 255
	var l [256]uint8
	for i := range l {
		a := float64(i) / 255
		dark := 1 - math.Pow(1-a, darkTextGamma)
		light := math.Pow(a, lightTextGamma)
		l[i] = uint8(math.Round((dark*(1-lum) + light*lum) * 255))
	}
	luts.Store(key, &l)
	return &l
}
