package properties

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// Mix blends fg over bg in sRGB space, as CSS and design tools do: t =
// 0 gives bg, t = 1 gives fg. Plain colors give an opaque
// "#RRGGBBAA" color. When either side is a theme token ("theme.surface")
// or another mix that uses one, Mix returns a CSS color-mix() value
// instead ("color-mix(in srgb, theme.text 8%, theme.surface)"), which
// resolves against the current palette each time it's painted - so a
// mixed color made once, at construction, follows light/dark switches
// like the tokens in it do.
//
// Use it instead of a translucent color for tints and hover shades.
// Gio composites translucent colors in linear light, which is
// physically correct but makes a faint overlay much brighter than in
// a browser: "#ffffff0d" (5% white) over "#101115" renders close to
// "#434344", where CSS shows about "#1c1d21". Mix("#ffffff",
// "#101115", 0.05) gives the CSS result.
func Mix(fg, bg Color, t float64) Color {
	t = max(0, min(1, t))
	if themed(fg) || themed(bg) {
		return fmt.Sprintf("color-mix(in srgb, %s %s%%, %s)", strings.TrimSpace(fg), strconv.FormatFloat(t*100, 'f', -1, 64), strings.TrimSpace(bg))
	}
	return CalcColorReverse(mixNRGBA(CalcColor(fg), CalcColor(bg), t))
}

// themed reports whether c depends on the palette.
func themed(c Color) bool {
	return strings.Contains(c, "theme.")
}

func mixNRGBA(a, b color.NRGBA, t float64) color.NRGBA {
	lerp := func(x, y uint8) uint8 {
		return uint8(float64(y) + (float64(x)-float64(y))*t + 0.5)
	}
	return color.NRGBA{
		R: lerp(a.R, b.R),
		G: lerp(a.G, b.G),
		B: lerp(a.B, b.B),
		A: lerp(a.A, b.A),
	}
}

// parseColorMix parses "color-mix(in srgb, A [p%], B [q%])", the colors
// being anything CalcColor takes. As in CSS, a missing percentage is
// 100% less the other one, and both missing mix half and half.
func parseColorMix(s string) (color.NRGBA, bool) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(strings.ToLower(s), "color-mix(") || !strings.HasSuffix(s, ")") {
		return color.NRGBA{}, false
	}
	args := splitTopLevel(s[len("color-mix(") : len(s)-1])
	if len(args) != 3 || !strings.EqualFold(strings.Join(strings.Fields(args[0]), " "), "in srgb") {
		return color.NRGBA{}, false
	}
	fg, p, okA := colorStop(args[1])
	bg, q, okB := colorStop(args[2])
	if !okA || !okB {
		return color.NRGBA{}, false
	}
	t := 0.5
	switch {
	case p >= 0:
		t = p
	case q >= 0:
		t = 1 - q
	}
	return mixNRGBA(CalcColor(fg), CalcColor(bg), max(0, min(1, t))), true
}

// colorStop splits "color [p%]" into the color and p/100, or -1 when
// there's no percentage.
func colorStop(s string) (Color, float64, bool) {
	s = strings.TrimSpace(s)
	i := strings.LastIndexByte(s, ' ')
	if i < 0 || !strings.HasSuffix(s, "%") || strings.HasSuffix(s[:i], ",") {
		return s, -1, s != ""
	}
	p, err := strconv.ParseFloat(s[i+1:len(s)-1], 64)
	if err != nil {
		return s, -1, true
	}
	return strings.TrimSpace(s[:i]), p / 100, true
}

// splitTopLevel splits s at commas outside parentheses.
func splitTopLevel(s string) []string {
	var out []string
	depth, start := 0, 0
	for i, r := range s {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				out = append(out, s[start:i])
				start = i + 1
			}
		}
	}
	return append(out, s[start:])
}
