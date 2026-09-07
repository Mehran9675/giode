package properties

import "gioui.org/font"

// FontWeight is the weight of text (CSS font-weight).
type FontWeight string

const (
	FontWeightThin      FontWeight = "thin"
	FontWeightLight     FontWeight = "light"
	FontWeightNormal    FontWeight = "normal"
	FontWeightMedium    FontWeight = "medium"
	FontWeightBold      FontWeight = "bold"
	FontWeightExtraBold FontWeight = "extrabold"
	FontWeightBlack     FontWeight = "black"
)

// Weight maps the value to a font weight, defaulting to normal.
func (w FontWeight) Weight() font.Weight {
	switch w {
	case FontWeightThin:
		return font.Thin
	case FontWeightLight:
		return font.Light
	case FontWeightMedium:
		return font.Medium
	case FontWeightBold:
		return font.Bold
	case FontWeightExtraBold:
		return font.ExtraBold
	case FontWeightBlack:
		return font.Black
	default:
		return font.Normal
	}
}
