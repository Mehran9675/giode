package main

import "github.com/mehran9675/giode"

// SpinnerShowcase is stateless: Spinner animates itself, so it can be
// built fresh every frame like Text.
func SpinnerShowcase() giode.Element {
	return giode.Stack(boxStyle,
		giode.Text("Spinner", textStyle),
		giode.Spinner(giode.Styles{Color: "#3b82f6", Width: 32, Height: 32}),
	)
}
