package kit

import "github.com/mehran9675/giode/styles/properties"

// Hovered returns the hover highlight of c: slightly lighter.
func Hovered(c properties.Color) properties.Color {
	return Scale(c, 1.15)
}
