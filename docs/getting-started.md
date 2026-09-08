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
	btn := giode.Button("Click me", giode.Styles{
		Background:   giode.HexColor("#3b82f6"),
		Color:        giode.HexColor("#ffffff"),
		BorderRadius: 8,
	}).OnClick(func() { count++ })

	app := giode.New(giode.Config{
		Title:      "Hello",
		Size:       image.Pt(800, 600),
		Resizable:  true,
		Background: giode.HexColor("#0f172a"),
	})

	if err := app.Run(func() giode.Element {
		return giode.Box(
			giode.Styles{Align: giode.AlignCenter, Justify: giode.JustifyCenter, Width: -1, Height: -1},
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

- **Stateless elements** (`Box`, `Row`, `Stack`, `Text`, `Divider`, `Responsive`, `Image`, `Progress`,
  `Spinner`, icons) are cheap to build inside the view function every frame.
- **Stateful components** (`Button`, `Tabs`, `Drawer`, `Dialog`, `Scroll`, `Menu`, `Router`,
  and the form controls under `Input`: `Input.Text`, `Input.Checkbox`, `Input.Slider`,
  `Input.Select`) hold their persistent state. Create them **once, before `Run`**, and lay
  them out every frame:

```go
input := giode.Input.Text("Name")    // once
...
app.Run(func() giode.Element {       // every frame
	return giode.Box(giode.Styles{}, input)
})
```

Creating a stateful component inside the view function resets it every frame.

## Examples

The repository's `examples/` directory has three runnable programs (`go run ./examples/<name>`):

- **hello** — the smallest possible app: a button, a counter and a right-click menu.
- **custom** — building a custom stateful component (a star rating) from scratch; see
  [Extending the library](extending.md).
- **showcase** — every built-in component in one window, each in its own file
  (`examples/showcase/*.go`) following the "construct once, `View()` every frame" pattern
  used throughout this guide. It also doubles as a live demo of `Wrap` (the gallery itself
  wraps as the window narrows) paired with `Scroll` for the resulting vertical overflow.
