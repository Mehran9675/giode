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
	btn := giode.Button("Click me",
		giode.Styles{
			Background:   "#3b82f6",
			Color:        "#ffffff",
			BorderRadius: 8,
		},
	).OnClick(func() { count++ })

	app := giode.New(giode.Config{
		Title:      "Hello",
		Size:       image.Pt(800, 600),
		Resizable:  true,
		Background: "#0f172a",
	})

	if err := app.Run(func() giode.Element {
		return giode.Stack(giode.Styles{
			Align:   giode.AlignCenter,
			Justify: giode.JustifyCenter,
			Width:   -1,
			Height:  -1,
		},
			giode.Text(fmt.Sprintf("Clicks: %d", count), giode.Styles{Color: "#f8fafc"}),
			btn,
		)
	}); err != nil {
		panic(err)
	}
}
```

## The two kinds of elements

Immediate mode rebuilds the UI every frame. Giode handles the two consequences explicitly:

- **Stateless elements** (`Box`, `Row`, `Stack`, `Text`, `Divider`, `Responsive`, `Image`,
  `Progress`, `Spinner`, icons) are cheap to build inside the view function every frame.
- **Stateful components** (`Button`, `Input.Text`, `Input.Checkbox`, `Input.Slider`,
  `Input.Select`, `Tabs`, `Scroll`, `Drawer`, `Dialog`, `Menu`, `Router`) hold their
  persistent state. Create them **once, before `Run`**, and lay them out every frame:

```go
name := giode.Input.Text("Name")  // once
...
app.Run(func() giode.Element {    // every frame
	return giode.Stack(giode.Styles{}, name)
})
```

Creating a stateful component inside the view function resets it every frame.

## Lifecycle

`Run` blocks for the lifetime of the app and terminates the process when the window
closes. Hook into the lifecycle with `Config.OnStartup`, `OnReady`, `OnShutdown` and
`OnBeforeClose`; see [Window configuration](window.md).
