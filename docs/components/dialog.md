# Dialog

Centered modal dialog with a dimmed scrim. It implements `Element`, so it drops
directly into the tree like any other element — no raw Gio required.

Visibility is bound to a `*bool` you own (external state): set it to open the dialog
from anywhere (a button, a menu item, ...), and the dialog itself clears it when the
scrim is clicked.

```go
dialogOpen := false // owned by you, alongside your other app state

dlg := giode.Dialog(&dialogOpen, func() giode.Element {
	return giode.Text("Are you sure?")
}, giode.Styles{})

confirm := giode.Button("Delete", giode.Styles{}).OnClick(func() { dialogOpen = true })

// in the view, dlg last so it renders above the UI (a no-op while closed):
giode.Stack(giode.Styles{Width: -1, Height: -1},
	mainContent,
	confirm,
	dlg,
)
```

## API

| Member | Description |
| --- | --- |
| `Dialog(open *bool, content func() Element, st ...Styles) *Dialog` | Creates the dialog, visible whenever `*open` is true. `content` is rebuilt every frame it is shown. The styles argument is optional and applies to the panel. |
| `Opened() bool` | Reports whether the dialog is currently visible. |

Clicking the scrim sets `*open` back to `false`. Panel styles follow `Box`: `Background`,
`BorderRadius` (default 12), `Padding` (default 20).
