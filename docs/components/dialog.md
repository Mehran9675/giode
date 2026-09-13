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

Clicking the scrim sets `*open` to false. Panel styles follow `Box`: `Background`,
`BorderRadius` (default 12), `Padding` (default 20).
