# Drawer

A panel that slides in from the left or right edge over a dimmed background. You own the
open state as a `*bool`.

```go
open := new(bool)
nav := giode.Drawer(open, giode.SideLeft, 260, func() giode.Element {
	return giode.Stack(giode.Styles{Padding: giode.UniformInset(16), Gap: 8}, links...)
})

menuButton := giode.Button("Menu").OnClick(func() { *open = true })

// in the view, after the page so it draws on top:
giode.Stack(giode.Styles{Width: -1, Height: -1}, page, nav)
```

| Method | Description |
| --- | --- |
| `Drawer(open, side, width, content, st...) *Drawer` | Creates the drawer. `side` is `SideLeft` or `SideRight`; `width` is in dp. |
| `Opened() bool` | Whether it's open. |

Opening and closing animate over 200ms. Clicking the dimmed background sets `*open` to
false. Styles apply to the panel; `Background` defaults to `#1e293b`.
