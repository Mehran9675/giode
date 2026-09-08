package main

import "github.com/mehran9675/giode"

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
		),
	)
}
