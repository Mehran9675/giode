# Dialog

A modal dialog centered over a dimmed background. You own the open state as a `*bool`, so
any code can open or close it.

```go
open := new(bool)
confirm := giode.Dialog(open, func() giode.Element {
	return giode.Stack(giode.Styles{Gap: 12},
		giode.Text("Delete this file?", giode.Styles{Color: "white"}),
		giode.Row(giode.Styles{Gap: 8}, cancelButton, deleteButton),
	)
})

*open = true // show it

// in the view, after the page so it draws on top:
giode.Stack(giode.Styles{Width: -1, Height: -1}, page, confirm)
```

| Method | Description |
| --- | --- |
| `Dialog(open, content, st...) *Dialog` | Creates the dialog. |
| `Opened() bool` | Whether it's open. |

Clicking the dimmed background sets `*open` to false; clicks inside the dialog don't.

Styles apply to the panel: `Background` (`#1e293b` by default), `BorderRadius` (12),
`Padding` (20), and the other box styles.
