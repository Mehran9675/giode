# Checkbox

Toggle with an optional label. Shows a check icon when checked.

```go
notify := giode.Input.Checkbox("Enable notifications").
	OnChange(func(checked bool) { ... })
notify.SetChecked(true)

// bare checkbox, no label, fully custom look:
star := giode.Input.Checkbox("", giode.Styles{
	Width: 24, Height: 24, Background: "green", BorderRadius: 6,
}).Icon("star", giode.Styles{Color: "yellow"})

// in the view:
giode.Box(giode.Styles{}, notify)
```

## API

| Member | Description |
| --- | --- |
| `Input.Checkbox(label string, st ...Styles) *Checkbox` | Creates the checkbox. `label` is optional: pass `""` for none. The styles argument is optional. |
| `Icon(name string, st ...Styles) *Checkbox` | Replaces the material icon shown when checked (default `"check"`) and its styles. Returns the checkbox for chaining. |
| `Checked() bool` | Current state. |
| `SetChecked(v bool)` | Sets the state without invoking OnChange. |
| `OnChange(fn func(checked bool)) *Checkbox` | Change handler. |

## Styles

| Property | Role |
| --- | --- |
| `Color` | Label color. |
| `Background` | Checked box fill. |
| `BorderColor` | Unchecked box border. |
| `BorderWidth` | Unchecked box border width (default 2). |
| `BorderRadius` | Box corner radius (default 4). |
| `Width` / `Height` | Box size (default 18x18). |
