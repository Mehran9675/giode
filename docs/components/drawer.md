# Drawer

Slide-out side panel with a dimmed scrim, animated over 200ms.

```go
drawer := giode.NewDrawer(giode.SideLeft, 260, giode.Styles{Background: giode.HexColor("#1e293b")})

// in the view, last so it renders above the UI:
giode.Raw(func(gtx layout.Context) layout.Dimensions {
	return drawer.Layout(gtx, giode.Box(
		giode.Styles{Padding: giode.UniformInset(16)},
		giode.Text("Drawer content"),
	))
})
```

## API

| Member | Description |
| --- | --- |
| `NewDrawer(side Side, width int, st Styles) *Drawer` | Creates the drawer; `side` is `SideLeft` or `SideRight`. |
| `Open()` / `Close()` / `Toggle()` | State control. |
| `Opened() bool` | State. |
| `Layout(gtx, content Element)` | Renders scrim and panel; content is rebuilt every frame while visible. |

Clicking the scrim closes the drawer.
