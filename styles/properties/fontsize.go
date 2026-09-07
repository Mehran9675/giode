package properties

import "gioui.org/unit"

// FontSize is the text size in sp (CSS font-size). Zero resolves to
// the default of 14sp.
type FontSize = unit.Sp

const defaultFontSize = unit.Sp(14)

// ResolveFontSize returns the effective font size.
func ResolveFontSize(s FontSize) unit.Sp {
	if s <= 0 {
		return defaultFontSize
	}
	return s
}
