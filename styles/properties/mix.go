package properties

import "image/color"

// Mix blends fg over bg in sRGB space, as CSS and design tools do,
// and returns the result as an opaque "#RRGGBBAA" color: t = 0 gives
// bg, t = 1 gives fg. Theme tokens are resolved when Mix is called, so
// call it from the view function to follow theme changes.
//
// Use it instead of a translucent color for tints and hover shades.
// Gio composites translucent colors in linear light, which is
// physically correct but makes a faint overlay much brighter than in
// a browser: "#ffffff0d" (5% white) over "#101115" renders close to
// "#434344", where CSS shows about "#1c1d21". Mix("#ffffff",
// "#101115", 0.05) gives the CSS result.
func Mix(fg, bg Color, t float64) Color {
	t = max(0, min(1, t))
	a, b := CalcColor(fg), CalcColor(bg)
	lerp := func(x, y uint8) uint8 {
		return uint8(float64(y) + (float64(x)-float64(y))*t + 0.5)
	}
	return CalcColorReverse(color.NRGBA{
		R: lerp(a.R, b.R),
		G: lerp(a.G, b.G),
		B: lerp(a.B, b.B),
		A: lerp(a.A, b.A),
	})
}
