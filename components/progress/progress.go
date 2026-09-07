// Package progress provides a determinate progress bar.
package progress

import (
	"image/color"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
)

var (
	defaultTrack = color.NRGBA{R: 0x33, G: 0x41, B: 0x55, A: 0xff}
	defaultFill  = color.NRGBA{R: 0x3b, G: 0x82, B: 0xf6, A: 0xff}
)

// Progress is a stateless determinate progress bar. It expands to the
// available width; the height defaults to 8px. Color is the fill and
// Background is the track.
func Progress(value float32, st ...styles.Styles) elements.Element {
	var s styles.Styles
	if len(st) > 0 {
		s = st[0]
	}
	return &progressEl{value: value, st: s}
}
