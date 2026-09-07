# Checkbox

Toggle with a label.

```go
notify := giode.NewCheckbox("Enable notifications").
	OnChange(func(checked bool) { ... })
notify.SetChecked(true)

// in the view:
giode.Box(giode.Styles{}, notify)
```

## API

| Member | Description |
| --- | --- |
| `NewCheckbox(label string) *Checkbox` | Creates the checkbox. |
| `Checked() bool` | Current state. |
| `SetChecked(v bool)` | Sets the state without invoking OnChange. |
| `OnChange(fn func(checked bool)) *Checkbox` | Change handler. |
| `Styles(st Styles) *Checkbox` | Styles. |

## Styles

| Property | Role |
| --- | --- |
| `Color` | Label color. |
| `Background` | Checked box fill. |
| `BorderColor` | Unchecked box border. |
