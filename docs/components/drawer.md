# Drawer

Slide-out side panel with a dimmed scrim, animated over 200ms. It implements
`Element`, so it drops directly into the tree like any other element — no raw Gio
required.

Visibility is bound to a `*bool` you own (external state): set it to open the drawer
from anywhere, and the drawer itself clears it when the scrim is clicked.

```go
drawerOpen := false // owned by you, alongside your other app state

drawer := giode.Drawer(&drawerOpen, giode.SideLeft, 260, func() giode.Element {
	return giode.Text("Drawer content")
}, giode.Styles{Background: giode.HexColor("#1e293b")})

open := giode.Button("Menu", giode.Styles{}).OnClick(func() { drawerOpen = true })

// in the view, drawer last so it renders above the UI (a no-op while closed):
giode.Stack(giode.Styles{Width: -1, Height: -1},
	mainContent,
	open,
	drawer,
)
```

## API

| Member | Description |
| --- | --- |
| `Drawer(open *bool, side Side, width int, content func() Element, st ...Styles) *Drawer` | Creates the drawer, open whenever `*open` is true; `side` is `SideLeft` or `SideRight`. `content` is rebuilt every frame it is shown. The styles argument is optional. |
| `Opened() bool` | Reports whether `*open` is currently true. |

Clicking the scrim sets `*open` back to `false`.
