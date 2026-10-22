# Tabs

A tab bar with one view per tab. Only the selected tab's view is built.

```go
tabs := giode.Tabs(giode.Styles{},
	giode.Tab{Label: "General", View: generalView},
	giode.Tab{Label: "Advanced", View: advancedView},
).OnChange(func(i int) { log.Println("tab", i) })
```

`View` is a `func() giode.Element`.

| Method | Description |
| --- | --- |
| `Tabs(st, tabs...) *Tabs` | Creates the tab bar. |
| `Selected() int` | Index of the selected tab. |
| `SetSelected(i)` | Selects a tab without calling `OnChange`. |
| `OnChange(fn)` | Called when the user picks a tab. |

Styles: `Color` (the underline of the selected tab), `BorderColor` (the selected tab's
label).
