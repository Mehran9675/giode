# Getting started

## Install

```sh
go get github.com/mehran9675/giode
```

Giode needs Go 1.26 or newer. It runs wherever Gio does: Windows, macOS and Linux, with
mobile support through Gio. On Linux, Gio needs a few system packages (Wayland, X11 and
Vulkan development headers); see [Gio's install notes](https://gioui.org/doc/install).

## A first app

```go
package main

import (
	"fmt"
	"image"
	"log"

	"github.com/mehran9675/giode"
)

func main() {
	giode.UseGofont() // load the Go fonts for text

	count := 0
	btn := giode.Button("Click me").OnClick(func() { count++ })

	app := giode.New(giode.Config{
		Title:      "Hello",
		Size:       image.Pt(800, 600),
		Resizable:  true,
		Background: "theme.background",
	})

	err := app.Run(func() giode.Element {
		return giode.Stack(giode.Styles{
			Width: -1, Height: -1, // fill the window
			Align: giode.AlignCenter, Justify: giode.JustifyCenter, Gap: 12,
		},
			giode.Text(fmt.Sprintf("Clicks: %d", count), giode.Styles{Color: "theme.text"}),
			btn,
		)
	})
	if err != nil {
		log.Fatal(err)
	}
}
```

The function passed to `Run` is the **view**. It runs every time the window needs a new
frame (on input, resize or a state change) and returns the whole UI.

## Elements and components

Because the view is rebuilt every frame, Giode has two kinds of building blocks:

- **Elements** hold no state: `Box`, `Row`, `Stack`, `Text`, `Divider`, `Image`, `Progress`,
  icons. Create them inside the view, every frame.
- **Components** remember things between frames: `Button` remembers whether it's pressed, a
  text field its text, `Scroll` its offset. Create them **once**, outside the view, and lay
  them out every frame:

```go
name := giode.Input.Text("Your name") // once

app.Run(func() giode.Element {
	return giode.Stack(giode.Styles{}, name) // every frame
})
```

A component created inside the view is a new one every frame, so it forgets everything:
clicks don't register and typed text disappears. The same goes for an element with hover
styles (`Box(...).OnHover(...)`), which keeps its hover state.

For lists whose items come and go, create per-item components with
[`Keyed`](state.md#per-item-state-for-dynamic-lists).

## Next steps

- [Elements](elements.md) and [Style properties](styles.md) for layout and styling.
- [Components](components/README.md) for the available widgets.
- [Window configuration](window.md) for size limits, frameless windows and lifecycle
  callbacks.
- [Threading](threading.md) before you load data in the background.
