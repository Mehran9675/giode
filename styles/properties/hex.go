package properties

import (
	"image/color"
	"strconv"
	"strings"
)

// HexColor validates s as a hex color and returns it unchanged as a
// Color. It panics on invalid input, same as Hex.
func HexColor(s string) Color {
	Hex(s)
	return s
}

// Hex parses a CSS hex color: "#RGB", "#RGBA", "#RRGGBB" or
// "#RRGGBBAA". The leading '#' is optional. It panics on invalid
// input.
func Hex(s string) color.NRGBA {
	h := strings.TrimPrefix(strings.TrimSpace(s), "#")
	switch len(h) {
	case 3, 4:
		b := make([]byte, 0, len(h)*2)
		for _, c := range h {
			b = append(b, byte(c), byte(c))
		}
		h = string(b)
	case 6, 8:
	default:
		panic("giode: invalid hex color " + strconv.Quote(s))
	}
	parse := func(off int) uint8 {
		v, err := strconv.ParseUint(h[off:off+2], 16, 8)
		if err != nil {
			panic("giode: invalid hex color " + strconv.Quote(s))
		}
		return uint8(v)
	}
	c := color.NRGBA{R: parse(0), G: parse(2), B: parse(4), A: 0xff}
	if len(h) == 8 {
		c.A = parse(6)
	}
	return c
}
