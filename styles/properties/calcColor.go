package properties

import (
	"image/color"
	"strings"
	"sync"
	"sync/atomic"
)

// fallbackColor is what an unparseable color resolves to: white, so a
// typo is visible rather than silently invisible.
var fallbackColor = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

// colorCache memoizes parsed colors. Colors are resolved many times per
// frame (every background, border and text run), so parsing the same
// string over and over would dominate the cost of painting a static
// UI. The cache is bounded: once it grows past maxCachedColors (which
// only happens when an app synthesizes colors dynamically, e.g. an
// animated hex string), it is dropped and refilled.
var (
	colorCache      sync.Map // string -> color.NRGBA
	colorCacheCount atomic.Int64
)

const maxCachedColors = 4096

// CalcColor resolves a Color value to its paintable form. An empty
// string means unset and resolves to fully transparent, so every
// "is this set?" check across the library (backgrounds, borders,
// shadows, component defaults) behaves consistently; an unrecognized
// non-empty value falls back to white so a typo is visible rather than
// silently invisible. CalcColor never panics.
//
// Accepted syntaxes are the CSS ones: hex ("#rgb", "#rgba",
// "#rrggbb", "#rrggbbaa"), functional ("rgb(59 130 246)",
// "rgba(59, 130, 246, 0.5)", "hsl(217 91% 60%)", "hsla(...)"), every
// CSS named color ("rebeccapurple") and "transparent".
//
// A "theme.<field>" value (e.g. "theme.primary", "theme.background")
// resolves against the active Palette instead of a fixed color - see
// SetMode/SetPalette - so restyling the palette or switching mode
// updates every element styled with that token, without editing each
// one. A token may carry an alpha suffix: "theme.primary/0.5" or
// "theme.primary/50%".
//
// "color-mix(in srgb, A p%, B)" mixes two colors as CSS does, A and B
// being any of the above; see Mix.
func CalcColor(c string) color.NRGBA {
	if c == "" {
		return color.NRGBA{}
	}
	if strings.HasPrefix(c, "theme.") {
		return resolveThemeColor(c)
	}
	if themed(c) {
		// A color-mix() of theme tokens changes with the palette.
		if col, ok := parseColorMix(c); ok {
			return col
		}
		return fallbackColor
	}
	if v, ok := colorCache.Load(c); ok {
		return v.(color.NRGBA)
	}
	col, err := ParseColor(c)
	if err != nil {
		col = fallbackColor
	}
	if colorCacheCount.Add(1) > maxCachedColors {
		colorCache.Clear()
		colorCacheCount.Store(1)
	}
	colorCache.Store(c, col)
	return col
}

// resolveThemeColor resolves "theme.<field>[/alpha]".
func resolveThemeColor(c string) color.NRGBA {
	token, alpha := c, float32(1)
	if i := strings.IndexByte(c, '/'); i >= 0 {
		token = c[:i]
		a, ok := parseAlpha(strings.TrimSpace(c[i+1:]))
		if !ok {
			return fallbackColor
		}
		alpha = a
	}
	resolved, ok := themeToken(token)
	if !ok || themed(resolved) {
		return fallbackColor
	}
	col := CalcColor(resolved)
	if alpha < 1 {
		col.A = uint8(float32(col.A)*alpha + 0.5)
	}
	return col
}

// IsSet reports whether c is a non-empty color value. Unlike checking
// CalcColor(c).A, it treats an explicit "transparent" as set, so a
// style can override an inherited color with transparency.
func IsSet(c Color) bool {
	return c != ""
}
