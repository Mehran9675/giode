# Button

A text button, optionally with icons.

```go
save := giode.Button("Save").
	Icon(giode.Icon("save", giode.Styles{Width: 16})).
	OnClick(func() { doc.Save() })
```

| Method | Description |
| --- | --- |
| `Button(label, st...)` | Creates the button. |
| `OnClick(fn)` | Click handler. |
| `Label(s)` | Changes the text. |
| `Variant(v)` | `VariantContained` (default), `VariantOutlined` or `VariantGhost`. |
| `Icon(el)`, `TrailingIcon(el)` | An element before or after the text, usually an icon. |
| `Disabled(true)` | Dims the button and ignores input. |
| `Loading(true)` | Shows a spinner in place of the icon and ignores clicks. |
| `IsDisabled()`, `IsLoading()` | Current state. |
| `Styles(st)` | Replaces the styles. |
| `Click()` | Clicks it from code (ignored while disabled or loading). |

Icons that don't set their own color take the text color.

## Variants

```go
giode.Button("Save")                                                     // contained
giode.Button("Cancel").Variant(giode.VariantOutlined)                    // border only
giode.Button("Skip", giode.Styles{Color: "theme.muted"}).Variant(giode.VariantGhost) // text only
```

- **Contained** fills with `Background` (`theme.primary` by default) and draws the text in
  `theme.onprimary`. With `BackgroundGradient` it fills with the gradient instead.
- **Outlined** is transparent with a border in `BorderColor` (or `Background`) and fills
  lightly on hover.
- **Ghost** shows only the text until hovered.

Hover and press shades are derived from the colors you set, gradient included.

## Loading state

```go
upload := giode.Button("Upload")
upload.OnClick(func() {
	upload.Loading(true)
	go func() {
		err := send(file)
		app.Post(func() {
			upload.Loading(false)
			report(err)
		})
	}()
})
```

## Styles

| Property | Role |
| --- | --- |
| `Background`, `BackgroundGradient` | Fill (contained) or accent color (other variants). |
| `Color` | Text and icon color. |
| `FontSize` | Text size. |
| `Padding` | 16×8 by default. |
| `Gap` | Space between icon and text, 8 by default. |
| `BorderRadius`, `BorderWidth`, `BorderColor` | As for any box. |
