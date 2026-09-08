package properties

import (
	"image/color"
	"strings"
)

var defaultColors = map[string]string{
	"white":   "#ffffff",
	"black":   "#000000",
	"gray":    "#858585",
	"yellow":  "#ffff00",
	"red":     "#ff0000",
	"green":   "#00ff00",
	"blue":    "#0000ff",
	"magenta": "#ff00ff",
	"cyan":    "#00ffff",
	"purple":  "#800080",
}

// CalcColor resolves a Color value to its paintable form. An empty
// string means unset and resolves to fully transparent, so every
// "is this set?" check across the library (backgrounds, borders,
// shadows, component defaults) behaves consistently; an unrecognized
// non-empty name falls back to white so a typo is visible rather than
// silently invisible.
func CalcColor(c string) color.NRGBA {
	if c == "" {
		return color.NRGBA{}
	}
	if !strings.HasPrefix(c, "#") {
		if value, ok := defaultColors[c]; ok {
			return Hex(value)
		}
		return Hex(defaultColors["white"])
	}

	return Hex(c)
}
