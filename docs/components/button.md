# Button

Clickable button.

```go
btn := giode.Button("Click me",
	giode.Styles{Background: "#3b82f6", BorderRadius: 8},
).OnClick(func() { count++ })

// in the view:
giode.Stack(giode.Styles{}, btn)
```

## API

| Member | Description |
| --- | --- |
| `Button(label string, st ...Styles) *Button` | Creates the button. |
| `OnClick(fn func()) *Button` | Click handler. |
| `Label(label string) *Button` | Replaces the label. |
| `Click()` | Simulates a click. |

## Styles

| Property | Role |
| --- | --- |
| `Background` | Face color; pressed and hovered shades are derived automatically. |
| `Color` | Label color. |
| `FontSize` | Label size. |
| `Padding` | Face padding; defaults to 16x8. |
| `BorderRadius`, `BorderWidth`, `BorderColor`, `Opacity`, `Cursor` | As usual. |
