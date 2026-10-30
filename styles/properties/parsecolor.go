package properties

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
)

// ErrInvalidColor is wrapped by every error ParseColor returns.
var ErrInvalidColor = errors.New("giode: invalid color")

// ParseColor parses a CSS color value: hex, rgb()/rgba(),
// hsl()/hsla(), a named color or "transparent". Theme tokens are not
// accepted here, since their value depends on the active palette; use
// CalcColor to resolve those. Use ParseColor to validate colors coming
// from configuration or user input before styling with them.
func ParseColor(s string) (color.NRGBA, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	switch {
	case v == "":
		return color.NRGBA{}, fmt.Errorf("%w: empty", ErrInvalidColor)
	case v == "transparent":
		return color.NRGBA{}, nil
	case v[0] == '#':
		return parseHex(v[1:], s)
	case strings.HasPrefix(v, "rgb"):
		return parseRGBFunc(v, s)
	case strings.HasPrefix(v, "hsl"):
		return parseHSLFunc(v, s)
	case strings.HasPrefix(v, "color-mix("):
		if c, ok := parseColorMix(s); ok {
			return c, nil
		}
		return color.NRGBA{}, fmt.Errorf("%w %q", ErrInvalidColor, s)
	}
	if hex, ok := namedColors[v]; ok {
		return parseHex(hex[1:], s)
	}
	// A bare hex string without '#' ("3b82f6") is accepted for
	// compatibility with Hex.
	if c, err := parseHex(v, s); err == nil {
		return c, nil
	}
	return color.NRGBA{}, fmt.Errorf("%w %q", ErrInvalidColor, s)
}

// MustParseColor is ParseColor, panicking on error.
func MustParseColor(s string) color.NRGBA {
	c, err := ParseColor(s)
	if err != nil {
		panic(err)
	}
	return c
}

func parseHex(h, orig string) (color.NRGBA, error) {
	switch len(h) {
	case 3, 4:
		var b [8]byte
		for i := 0; i < len(h); i++ {
			b[2*i], b[2*i+1] = h[i], h[i]
		}
		h = string(b[:2*len(h)])
	case 6, 8:
	default:
		return color.NRGBA{}, fmt.Errorf("%w %q", ErrInvalidColor, orig)
	}
	var out [4]uint8
	out[3] = 0xff
	for i := 0; i < len(h)/2; i++ {
		hi, ok1 := hexNibble(h[2*i])
		lo, ok2 := hexNibble(h[2*i+1])
		if !ok1 || !ok2 {
			return color.NRGBA{}, fmt.Errorf("%w %q", ErrInvalidColor, orig)
		}
		out[i] = hi<<4 | lo
	}
	return color.NRGBA{R: out[0], G: out[1], B: out[2], A: out[3]}, nil
}

func hexNibble(c byte) (uint8, bool) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, true
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}

// funcArgs splits "name(a, b, c / d)" or "name(a b c / d)" into its
// arguments; the alpha after '/' (if any) is returned as the last
// element.
func funcArgs(v string) ([]string, bool) {
	open := strings.IndexByte(v, '(')
	if open < 0 || !strings.HasSuffix(v, ")") {
		return nil, false
	}
	body := v[open+1 : len(v)-1]
	body = strings.ReplaceAll(body, "/", " ")
	body = strings.ReplaceAll(body, ",", " ")
	return strings.Fields(body), true
}

func parseRGBFunc(v, orig string) (color.NRGBA, error) {
	args, ok := funcArgs(v)
	if !ok || (len(args) != 3 && len(args) != 4) {
		return color.NRGBA{}, fmt.Errorf("%w %q", ErrInvalidColor, orig)
	}
	var ch [3]uint8
	for i := 0; i < 3; i++ {
		f, ok := parseChannel(args[i])
		if !ok {
			return color.NRGBA{}, fmt.Errorf("%w %q", ErrInvalidColor, orig)
		}
		ch[i] = f
	}
	a := float32(1)
	if len(args) == 4 {
		if a, ok = parseAlpha(args[3]); !ok {
			return color.NRGBA{}, fmt.Errorf("%w %q", ErrInvalidColor, orig)
		}
	}
	return color.NRGBA{R: ch[0], G: ch[1], B: ch[2], A: uint8(a*255 + 0.5)}, nil
}

func parseHSLFunc(v, orig string) (color.NRGBA, error) {
	args, ok := funcArgs(v)
	if !ok || (len(args) != 3 && len(args) != 4) {
		return color.NRGBA{}, fmt.Errorf("%w %q", ErrInvalidColor, orig)
	}
	h, err := strconv.ParseFloat(strings.TrimSuffix(args[0], "deg"), 64)
	if err != nil {
		return color.NRGBA{}, fmt.Errorf("%w %q", ErrInvalidColor, orig)
	}
	s, ok1 := parsePercent(args[1])
	l, ok2 := parsePercent(args[2])
	if !ok1 || !ok2 {
		return color.NRGBA{}, fmt.Errorf("%w %q", ErrInvalidColor, orig)
	}
	a := float32(1)
	if len(args) == 4 {
		if a, ok = parseAlpha(args[3]); !ok {
			return color.NRGBA{}, fmt.Errorf("%w %q", ErrInvalidColor, orig)
		}
	}
	r, g, b := hslToRGB(h, s, l)
	return color.NRGBA{R: r, G: g, B: b, A: uint8(a*255 + 0.5)}, nil
}

// parseChannel parses an rgb() channel: 0-255 or a percentage.
func parseChannel(s string) (uint8, bool) {
	if strings.HasSuffix(s, "%") {
		p, ok := parsePercent(s)
		if !ok {
			return 0, false
		}
		return uint8(p*255 + 0.5), true
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return uint8(math.Max(0, math.Min(255, f)) + 0.5), true
}

// parsePercent parses "50%" (or a bare number) into the range 0-1.
func parsePercent(s string) (float64, bool) {
	f, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
	if err != nil {
		return 0, false
	}
	return math.Max(0, math.Min(1, f/100)), true
}

// parseAlpha parses an alpha value: a number in 0-1 or a percentage.
func parseAlpha(s string) (float32, bool) {
	if strings.HasSuffix(s, "%") {
		p, ok := parsePercent(s)
		return float32(p), ok
	}
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, false
	}
	return float32(math.Max(0, math.Min(1, f))), true
}

func hslToRGB(h, s, l float64) (r, g, b uint8) {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2
	var rf, gf, bf float64
	switch {
	case h < 60:
		rf, gf, bf = c, x, 0
	case h < 120:
		rf, gf, bf = x, c, 0
	case h < 180:
		rf, gf, bf = 0, c, x
	case h < 240:
		rf, gf, bf = 0, x, c
	case h < 300:
		rf, gf, bf = x, 0, c
	default:
		rf, gf, bf = c, 0, x
	}
	to := func(v float64) uint8 { return uint8(math.Round((v + m) * 255)) }
	return to(rf), to(gf), to(bf)
}

// namedColors is the CSS Color Module Level 4 named color set.
var namedColors = map[string]string{
	"aliceblue": "#f0f8ff", "antiquewhite": "#faebd7", "aqua": "#00ffff", "aquamarine": "#7fffd4",
	"azure": "#f0ffff", "beige": "#f5f5dc", "bisque": "#ffe4c4", "black": "#000000",
	"blanchedalmond": "#ffebcd", "blue": "#0000ff", "blueviolet": "#8a2be2", "brown": "#a52a2a",
	"burlywood": "#deb887", "cadetblue": "#5f9ea0", "chartreuse": "#7fff00", "chocolate": "#d2691e",
	"coral": "#ff7f50", "cornflowerblue": "#6495ed", "cornsilk": "#fff8dc", "crimson": "#dc143c",
	"cyan": "#00ffff", "darkblue": "#00008b", "darkcyan": "#008b8b", "darkgoldenrod": "#b8860b",
	"darkgray": "#a9a9a9", "darkgreen": "#006400", "darkgrey": "#a9a9a9", "darkkhaki": "#bdb76b",
	"darkmagenta": "#8b008b", "darkolivegreen": "#556b2f", "darkorange": "#ff8c00", "darkorchid": "#9932cc",
	"darkred": "#8b0000", "darksalmon": "#e9967a", "darkseagreen": "#8fbc8f", "darkslateblue": "#483d8b",
	"darkslategray": "#2f4f4f", "darkslategrey": "#2f4f4f", "darkturquoise": "#00ced1", "darkviolet": "#9400d3",
	"deeppink": "#ff1493", "deepskyblue": "#00bfff", "dimgray": "#696969", "dimgrey": "#696969",
	"dodgerblue": "#1e90ff", "firebrick": "#b22222", "floralwhite": "#fffaf0", "forestgreen": "#228b22",
	"fuchsia": "#ff00ff", "gainsboro": "#dcdcdc", "ghostwhite": "#f8f8ff", "gold": "#ffd700",
	"goldenrod": "#daa520", "gray": "#808080", "green": "#008000", "greenyellow": "#adff2f",
	"grey": "#808080", "honeydew": "#f0fff0", "hotpink": "#ff69b4", "indianred": "#cd5c5c",
	"indigo": "#4b0082", "ivory": "#fffff0", "khaki": "#f0e68c", "lavender": "#e6e6fa",
	"lavenderblush": "#fff0f5", "lawngreen": "#7cfc00", "lemonchiffon": "#fffacd", "lightblue": "#add8e6",
	"lightcoral": "#f08080", "lightcyan": "#e0ffff", "lightgoldenrodyellow": "#fafad2", "lightgray": "#d3d3d3",
	"lightgreen": "#90ee90", "lightgrey": "#d3d3d3", "lightpink": "#ffb6c1", "lightsalmon": "#ffa07a",
	"lightseagreen": "#20b2aa", "lightskyblue": "#87cefa", "lightslategray": "#778899", "lightslategrey": "#778899",
	"lightsteelblue": "#b0c4de", "lightyellow": "#ffffe0", "lime": "#00ff00", "limegreen": "#32cd32",
	"linen": "#faf0e6", "magenta": "#ff00ff", "maroon": "#800000", "mediumaquamarine": "#66cdaa",
	"mediumblue": "#0000cd", "mediumorchid": "#ba55d3", "mediumpurple": "#9370db", "mediumseagreen": "#3cb371",
	"mediumslateblue": "#7b68ee", "mediumspringgreen": "#00fa9a", "mediumturquoise": "#48d1cc", "mediumvioletred": "#c71585",
	"midnightblue": "#191970", "mintcream": "#f5fffa", "mistyrose": "#ffe4e1", "moccasin": "#ffe4b5",
	"navajowhite": "#ffdead", "navy": "#000080", "oldlace": "#fdf5e6", "olive": "#808000",
	"olivedrab": "#6b8e23", "orange": "#ffa500", "orangered": "#ff4500", "orchid": "#da70d6",
	"palegoldenrod": "#eee8aa", "palegreen": "#98fb98", "paleturquoise": "#afeeee", "palevioletred": "#db7093",
	"papayawhip": "#ffefd5", "peachpuff": "#ffdab9", "peru": "#cd853f", "pink": "#ffc0cb",
	"plum": "#dda0dd", "powderblue": "#b0e0e6", "purple": "#800080", "rebeccapurple": "#663399",
	"red": "#ff0000", "rosybrown": "#bc8f8f", "royalblue": "#4169e1", "saddlebrown": "#8b4513",
	"salmon": "#fa8072", "sandybrown": "#f4a460", "seagreen": "#2e8b57", "seashell": "#fff5ee",
	"sienna": "#a0522d", "silver": "#c0c0c0", "skyblue": "#87ceeb", "slateblue": "#6a5acd",
	"slategray": "#708090", "slategrey": "#708090", "snow": "#fffafa", "springgreen": "#00ff7f",
	"steelblue": "#4682b4", "tan": "#d2b48c", "teal": "#008080", "thistle": "#d8bfd8",
	"tomato": "#ff6347", "turquoise": "#40e0d0", "violet": "#ee82ee", "wheat": "#f5deb3",
	"white": "#ffffff", "whitesmoke": "#f5f5f5", "yellow": "#ffff00", "yellowgreen": "#9acd32",
}
