# Drawer

Slide-out side panel with a dimmed scrim, animated over 200ms. The open state lives in a
`*bool` you own, so opening and closing works from anywhere.

```go
open := new(bool)
drawer := giode.Drawer(open, giode.SideLeft, 260, func() giode.Element {
	return giode.Stack(giode.Styles{Padding: giode.UniformInset(16)},
		giode.Text("Drawer content"),
	)
}, giode.Styles{Background: "#1e293b"})

*open = true  // or false

// in the view, last so it renders above the UI:
giode.Stack(giode.Styles{}, page, drawer)
```

## API

| Member | Description |
| --- | --- |
| `Drawer(open *bool, side Side, width int, content func() Element, st ...Styles) *Drawer` | Creates the drawer; `side` is `SideLeft` or `SideRight`. |
| `Opened() bool` | Whether the drawer is open. |

Clicking the scrim sets `*open` to false.
