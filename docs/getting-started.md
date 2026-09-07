# Getting started

## Install

```sh
go get github.com/mehran9675/giode
```

Requires Go 1.22+. Windows, macOS, Linux and mobile are supported through Gio.

## Hello world

```go
package main

import (
	"fmt"
	"image"

	"github.com/mehran9675/giode"
)

func main() {
	giode.UseGofont() // embed a font for text

	count := 0
	btn := giode.NewButton("Click me").
		OnClick(func() { count++ }).
		Styles(giode.Styles{
			Background:   giode.HexColor("#3b82f6"),
			Color:        giode.HexColor("#ffffff"),
			BorderRadius: 8,
		})

	app := giode.New(giode.Config{
		Title:      "Hello",
		Size:       image.Pt(800, 600),
		Resizable:  true,
		Background: giode.HexColor("#0f172a"),
	})

	if err := app.Run(func() giode.Element {
		return giode.Box(
			giode.Styles{Display: giode.DisplayFlex, Align: giode.AlignCenter, Justify: giode.JustifyCenter, Width: -1, Height: -1},
			giode.Text(fmt.Sprintf("Clicks: %d", count), giode.Styles{Color: giode.HexColor("#f8fafc")}),
			btn,
		)
	}); err != nil {
		panic(err)
	}
}
```

## The two kinds of elements

Immediate mode rebuilds the UI every frame. Giode handles the two consequences explicitly:

- **Stateless elements** (`Box`, `Flex`, `Text`, `Divider`, `Responsive`, `Image`, `Progress`,
  `Spinner`, icons) are cheap to build inside the view function every frame.
- **Stateful components** (`Button`, `Input`, `Checkbox`, `Slider`, `Tabs`, `Dropdown`,
  `Drawer`, `Dialog`, `Scroll`, `Menu`, `Router`) hold their persistent state. Create them
  **once, before `Run`**, and lay them out every frame:

```go
input := giode.NewInput("Name")      // once
...
app.Run(func() giode.Element {       // every frame
	return giode.Box(giode.Styles{}, input)
})
```

Creating a stateful component inside the view function resets it every frame.
