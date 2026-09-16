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
| `Variant(v ButtonVariant) *Button` | Visual style: `VariantContained` (default), `VariantOutlined` or `VariantGhost`. |
| `Click()` | Simulates a click. |

## Variants

```go
giode.Button("Save", giode.Styles{Background: "#3b82f6", Color: "white"})                              // contained (default)
giode.Button("Cancel", giode.Styles{BorderColor: "#3b82f6", Color: "#3b82f6"}).Variant(giode.VariantOutlined)
giode.Button("Skip", giode.Styles{Background: "#3b82f6", Color: "white"}).Variant(giode.VariantGhost)
```

- `VariantContained` fills with `Background` (or black if unset) at all times - the
  original, and only, look before `Variant` existed.
- `VariantOutlined` is transparent with a visible border (`BorderColor`, or `Background`
  if unset; `BorderWidth` defaults to 1), filling faintly with the border color on
  hover/press.
- `VariantGhost` is fully transparent at rest - only the label is visible - filling
  faintly with `Background` on hover/press.

## Styles

| Property | Role |
| --- | --- |
| `Background` | Face color (`VariantContained`) or accent color (`VariantOutlined`/`VariantGhost`); pressed and hovered shades are derived automatically. |
| `Color` | Label color. |
| `FontSize` | Label size. |
| `Padding` | Face padding; defaults to 16x8. |
| `BorderRadius`, `BorderWidth`, `BorderColor`, `Opacity`, `Cursor` | As usual. |
