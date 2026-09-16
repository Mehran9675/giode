package properties

import "strings"

// Mode is a global color mode: which Palette "theme.*" color tokens
// (see CalcColor) resolve against.
type Mode string

const (
	// ModeLight is the default mode.
	ModeLight Mode = "light"
	ModeDark  Mode = "dark"
)

// Palette is a named set of semantic colors that "theme.*" color
// tokens (e.g. "theme.primary", "theme.background") resolve against,
// so a component styled once with a token automatically follows
// SetMode/SetPalette instead of a fixed color. Each field takes a
// plain color - hex or a named CSS color, like any other Color - not
// another token.
type Palette struct {
	Primary    Color
	Secondary  Color
	Background Color
	Surface    Color
	Text       Color
	Border     Color
}

// DefaultLightPalette and DefaultDarkPalette are the built-in
// palettes for each Mode, in effect until SetPalette overrides one.
var (
	DefaultLightPalette = Palette{
		Primary:    "#3b82f6",
		Secondary:  "#64748b",
		Background: "#ffffff",
		Surface:    "#f1f5f9",
		Text:       "#0f172a",
		Border:     "#cbd5e1",
	}
	DefaultDarkPalette = Palette{
		Primary:    "#3b82f6",
		Secondary:  "#94a3b8",
		Background: "#0f172a",
		Surface:    "#1e293b",
		Text:       "#f8fafc",
		Border:     "#334155",
	}
)

var (
	currentMode  = ModeLight
	lightPalette = DefaultLightPalette
	darkPalette  = DefaultDarkPalette
)

// SetMode switches the active color Mode, changing what every
// "theme.*" color token in the tree resolves to on the next frame.
func SetMode(m Mode) {
	currentMode = m
}

// CurrentMode returns the active color mode (ModeLight by default).
func CurrentMode() Mode {
	return currentMode
}

// SetPalette replaces the palette used for mode m, changing what
// every "theme.*" color token resolves to while that mode is active.
func SetPalette(m Mode, p Palette) {
	if m == ModeDark {
		darkPalette = p
	} else {
		lightPalette = p
	}
}

// currentPalette returns the palette for the active mode.
func currentPalette() Palette {
	if currentMode == ModeDark {
		return darkPalette
	}
	return lightPalette
}

// themeToken resolves a "theme.<field>" color token against the
// active palette's matching field, case-insensitively. ok is false
// for anything that isn't a recognized token.
func themeToken(c string) (Color, bool) {
	p := currentPalette()
	switch strings.ToLower(strings.TrimPrefix(c, "theme.")) {
	case "primary":
		return p.Primary, true
	case "secondary":
		return p.Secondary, true
	case "background":
		return p.Background, true
	case "surface":
		return p.Surface, true
	case "text":
		return p.Text, true
	case "border":
		return p.Border, true
	default:
		return "", false
	}
}
