# Slider

Horizontal slider in the range [0, 1]. Expands to the available width.

```go
volume := giode.NewSlider().
	OnChange(func(v float32) { player.SetVolume(v) })

// in the view:
giode.Box(giode.Styles{}, volume)
```

## API

| Member | Description |
| --- | --- |
| `NewSlider() *Slider` | Creates the slider. |
| `Value() float32` | Current value in [0, 1]. |
| `SetValue(v float32)` | Sets the value without invoking OnChange. |
| `OnChange(fn func(v float32)) *Slider` | Called while dragging. |
| `Styles(st Styles) *Slider` | Styles. |

## Styles

| Property | Role |
| --- | --- |
| `Color` | Fill up to the thumb. |
| `Background` | Track. |
| `BorderColor` | Thumb. |
