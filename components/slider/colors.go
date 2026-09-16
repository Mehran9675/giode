package slider

import "github.com/mehran9675/giode/styles/properties"

var (
	defaultTrack = "#334155"
	defaultFill  = "#3b82f6"
	defaultThumb = "#f8fafc"
)

const (
	height  = 24
	trackH  = 4
	thumbD  = 14
	padding = thumbD / 2
)

func resolve(c, def properties.Color) properties.Color {
	if c == "" {
		return def
	}
	return c
}
