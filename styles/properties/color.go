package properties

import "image/color"

// Color is the text color (CSS color). A zero alpha color resolves to
// black, the CSS default.
type Color = string

// ResolveColor returns the effective text color.
func ResolveColor(c Color) color.NRGBA {
	if CalcColor(c).A == 0 {
		return color.NRGBA{A: 0xff}
	}
	return CalcColor(c)
}
