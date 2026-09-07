# Dialog

Centered modal dialog with a dimmed scrim.

```go
dialog := giode.NewDialog(giode.Styles{})

// in the view, last so it renders above the UI:
giode.Raw(func(gtx layout.Context) layout.Dimensions {
	return dialog.Layout(gtx, giode.Box(
		giode.Styles{},
		giode.Text("Are you sure?"),
	))
})
```

## API

| Member | Description |
| --- | --- |
| `NewDialog(st Styles) *Dialog` | Creates the dialog; styles apply to the panel. |
| `Open()` / `Close()` | State control. |
| `Opened() bool` | State. |
| `Layout(gtx, content Element)` | Renders the scrim and centered panel. |

Clicking the scrim closes the dialog. Panel styles follow `Box`: `Background`,
`BorderRadius` (default 12), `Padding` (default 20).
