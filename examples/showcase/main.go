package main

import (
	"image"
	"log"

	"github.com/mehran9675/giode"
)

var boxStyle = giode.Styles{
	Align:        giode.AlignEnd,
	Justify:      giode.JustifyCenter,
	Padding:      giode.SymmetricInset(10, 10),
	BoxShadow:    giode.BoxShadow{X: 1, Y: 1, Blur: 4, Color: "black"},
	Background:   "black",
	BorderRadius: 5,
	Gap:          10,
	Width:        200,
}
var textStyle = giode.Styles{Color: "white", FontSize: 24, TextAlign: giode.TextAlignCenter}

func main() {
	giode.UseGofont()

	app := giode.New(giode.Config{
		Title:     "GIODE SHOWCASE",
		Size:      image.Pt(800, 600),
		Resizable: true,
		Frameless: false,
		Draggable: true,
	})

	buttonPg := NewButtonPage()
	checkboxPg := NewCheckboxPage()
	inputPg := NewInputPage()
	sliderPg := NewSliderPage()
	progressPg := NewProgressPage(sliderPg) // reads sliderPg's live value
	dialogPg := NewDialogShowcase()
	dropdownPg := NewDropdownPage()
	drawerPg := NewDrawerPage()
	tabsPg := NewTabsPage()
	scrollPg := NewScrollPage()
	routerPg := NewRouterPage()

	app.SetContextMenu(NewShowcaseMenu())

	gallery := giode.Scroll(func() giode.Element {
		return giode.Row(giode.Styles{
			Width:   -1,
			Gap:     20,
			Padding: giode.SymmetricInset(10, 10),
		},
			TextShowcase(),
			buttonPg.View(),
			checkboxPg.View(),
			inputPg.View(),
			sliderPg.View(),
			progressPg.View(),
			dialogPg.View(),
			dropdownPg.View(),
			drawerPg.View(),
			tabsPg.View(),
			scrollPg.View(),
			SpinnerShowcase(),
			ImageShowcase(),
			IconsShowcase(),
			DividerShowcase(),
			ResponsiveShowcase(),
			routerPg.View(),
			MenuHint(),
		)
	}, giode.Styles{ScrollBar: giode.ScrollBar{ThumbColor: "#3b82f6"}})

	if err := app.Run(func() giode.Element {
		return giode.Stack(giode.Styles{Width: -1, Height: -1},
			giode.Box(giode.Styles{Width: -1, Height: -1, Background: "#0f172a"}, gallery),
			drawerPg.Overlay(), // last, so overlays paint above the rest
			dialogPg.Overlay(),
		)
	}); err != nil {
		log.Fatal(err)
	}
}
