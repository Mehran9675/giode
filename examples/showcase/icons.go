package main

import "github.com/mehran9675/giode"

// docIcon is a full .svg document (viewBox + two separate <path>
// elements) loaded with LoadIconSVG, unlike the other icons here
// which are built from raw path data.
var docIcon = giode.MustLoadIconSVG([]byte(`<svg viewBox="0 0 24 24">
	<path d="M2 2 L10 2 L10 10 L2 10 Z"/>
	<path d="M14 14 L22 14 L22 22 L14 22 Z"/>
</svg>`))

// IconsShowcase is stateless: icons are stateless elements, so it can
// be built fresh every frame like Text.
func IconsShowcase() giode.Element {
	return giode.Stack(boxStyle,
		giode.Text("Icons", textStyle),
		giode.Row(giode.Styles{Gap: 8},
			giode.IconMaterial("star", giode.Styles{Color: "#facc15", Width: 22, Height: 22}),
			giode.IconMaterial("heart", giode.Styles{Color: "#ef4444", Width: 22, Height: 22}),
			giode.IconMaterial("settings", giode.Styles{Color: "white", Width: 22, Height: 22}),
			giode.IconMaterial("trash", giode.Styles{Color: "#94a3b8", Width: 22, Height: 22}),
			docIcon.Icon(giode.Styles{Color: "#22d3ee", Width: 22, Height: 22}),
		),
	)
}
