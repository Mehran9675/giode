package main

import (
	"log"

	"github.com/mehran9675/giode"
	giodemenu "github.com/mehran9675/giode/components/menu"
)

// NewShowcaseMenu builds the global right-click context menu. Unlike
// the other components it is not placed in the view tree: register it
// once with App.SetContextMenu and the app renders it automatically.
func NewShowcaseMenu() *giodemenu.Menu {
	return giode.Menu().
		IconItem("star", "Starred action", func() { log.Println("starred") }).
		Item("Plain action", func() { log.Println("clicked") }).
		Separator().
		DisabledItem("Disabled")
}

func MenuHint() giode.Element {
	return giode.Stack(boxStyle,
		giode.Text("Menu", textStyle),
		giode.Text("Right-click anywhere", giode.Styles{Color: "white"}),
	)
}
