# Dialog

Centered modal dialog with a dimmed scrim. The open state lives in a `*bool` you own.

```go
open := new(bool)
dialog := giode.Dialog(open, func() giode.Element {
	return giode.Stack(giode.Styles{},
		giode.Text("Are you sure?"),
	)
}, giode.Styles{})

*open = true  // or false

// in the view, last so it renders above the UI:
giode.Stack(giode.Styles{}, page, dialog)
```

## API

| Member | Description |
| --- | --- |
| `Dialog(open *bool, content func() Element, st ...Styles) *Dialog` | Creates the dialog; styles apply to the panel. |
| `Opened() bool` | Whether the dialog is open. |
| `Dismissible(on bool) *Dialog` | Whether the scrim, Back and Escape close it (on by default). |

Clicking the scrim, pressing Escape or pressing the system Back button (Android) sets
`*open` to false. Turn this off with `Dismissible(false)` for a dialog the user must
answer through its own buttons. Back and Escape go to the topmost open overlay only. Panel styles follow `Box`: `Background`,
`BorderRadius` (default 12), `Padding` (default 20).
