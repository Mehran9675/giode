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
// input; use ParseColor for a non-panicking, error-returning parser
// that accepts every CSS color syntax.
func Hex(s string) color.NRGBA {
	c, err := parseHex(strings.TrimPrefix(strings.TrimSpace(s), "#"), s)
	if err != nil {
		panic("giode: invalid hex color " + strconv.Quote(s))
	}
	return c
}
