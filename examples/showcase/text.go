package main

import "github.com/mehran9675/giode"

func TextShowcase() giode.Element {
	return giode.Stack(
		boxStyle,
		giode.Text("This is a Text", textStyle),
		giode.H1("This is a H1", textStyle),
		giode.H2("This is a H2", textStyle),
		giode.H3("This is a H3", textStyle),
		giode.H4("This is a H4", textStyle),
		giode.H5("This is a H5", textStyle),
		giode.H6("This is a H6", textStyle),
	)
}
