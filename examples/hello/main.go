// Command hello is a minimal Giode demo: a frameless, draggable
// window with a counter button and a global right-click menu.
package main

import (
	"fmt"
	"image"
	"log"

	"github.com/mehran9675/giode"
)

func main() {
	giode.UseGofont()

	count := 0
	btn := giode.NewButton("Click me").
		OnClick(func() { count++ }).
		Styles(giode.Styles{
			Background:   giode.HexColor("#3b82f6"),
			Color:        giode.HexColor("#ffffff"),
			Padding:      giode.SymmetricInset(16, 8),
			BorderRadius: 8,
		})

	menu := giode.NewMenu().
		Item("Reset counter", func() { count = 0 }).
		Separator().
		Item("Quit", func() { log.Println("bye") })

	app := giode.New(giode.Config{
		Title:      "Giode",
		Size:       image.Pt(800, 600),
		Resizable:  true,
		Frameless:  true,
		Draggable:  true,
		Background: giode.HexColor("#0f172a"),
	})
	app.SetContextMenu(menu)

	if err := app.Run(func() giode.Element {
		return giode.Box(
			giode.Styles{
				Display: giode.DisplayFlex,
				Align:   giode.AlignCenter,
				Justify: giode.JustifyCenter,
				Width:   -1,
				Height:  -1,
			},
			giode.Box(
				giode.Styles{
					Display:      giode.DisplayFlex,
					Direction:    giode.DirColumn,
					Align:        giode.AlignCenter,
					Gap:          12,
					Padding:      giode.UniformInset(24),
					Background:   giode.HexColor("#1e293b"),
					BorderRadius: 12,
					BoxShadow:    giode.BoxShadow{Y: 4, Blur: 12, Color: giode.HexColor("#00000060")},
				},
				giode.Text("Welcome to Giode", giode.Styles{Color: giode.HexColor("#f8fafc"), FontSize: 24}),
				giode.Text(fmt.Sprintf("Clicks: %d", count), giode.Styles{Color: giode.HexColor("#94a3b8")}),
				btn,
			),
		)
	}); err != nil {
		log.Fatal(err)
	}
}
