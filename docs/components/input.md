# Input

Form controls are grouped under the `Input` API:

| Constructor | Component |
| --- | --- |
| `giode.Input.Text(placeholder, st...)` | Single-line text field |
| `giode.Input.Checkbox(label, st...)` | Checkbox (label may be `""`) |
| `giode.Input.Slider(st...)` | Horizontal slider in [0, 1] |
| `giode.Input.Select(st, options...)` | Dropdown among options |

## Text field

```go
name := giode.Input.Text("Your name").
	Submit(func(text string) { save(text) })

// in the view:
giode.Stack(giode.Styles{}, name)
```

| Member | Description |
| --- | --- |
| `Text() string` | Current content. |
| `SetText(s string)` | Replaces the content. |
| `Submit(fn func(text string))` | Handler for Enter. |
| `Focus()` | Requests keyboard focus. |

Styles: `Color` (text), `Background` (fill), `BorderWidth`/`BorderColor`/`BorderRadius`
(frame, 1px border by default), `Padding` (12x8 by default). Expands to the available
width; set `Width` to fix it.

## Checkbox

```go
notify := giode.Input.Checkbox("Enable notifications").
	OnChange(func(checked bool) { ... })
notify.SetChecked(true)
```

| Member | Description |
| --- | --- |
| `Checked() bool` | Current state. |
| `SetChecked(v bool)` | Sets the state without invoking OnChange. |
| `OnChange(fn func(bool))` | Change handler. |
| `Icon(name string, st ...Styles)` | Replaces the check mark with a material icon. |

Styles: `Color` (label), `Background` (checked fill), `BorderColor` (unchecked border).

## Slider

```go
volume := giode.Input.Slider().
	OnChange(func(v float32) { player.SetVolume(v) })
```

| Member | Description |
| --- | --- |
| `Value() float32` | Current value in [0, 1]. |
| `SetValue(v float32)` | Sets the value without invoking OnChange. |
| `OnChange(fn func(float32))` | Called while dragging. |

Styles: `Color` (fill), `Background` (track), `BorderColor` (thumb).

## Select

```go
size := giode.Input.Select(giode.Styles{}, "Small", "Medium", "Large").
	OnChange(func(i int) { ... })
```

| Member | Description |
| --- | --- |
| `Selected() int` | Selected index. |
| `Value() string` | Selected option text. |
| `SetSelected(i int)` | Changes the selection without invoking OnChange. |
| `OnChange(fn func(int))` | Change handler. |
| `Opened() bool` | Whether the panel is open. |

Styles: `Background` (field and panel fill), `Color` (text), `BorderColor` (field
border), `BorderRadius` (corners), `Width` (fixed field width). The panel opens below
the field, clamped to the bounds of the parent container; siblings laid out after the
select paint above an open panel.
