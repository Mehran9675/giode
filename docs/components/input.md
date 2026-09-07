# Input

Single-line text field.

```go
name := giode.NewInput("Your name").
	Submit(func(text string) { save(text) })

// in the view:
giode.Box(giode.Styles{}, name)
```

## API

| Member | Description |
| --- | --- |
| `NewInput(placeholder string) *Input` | Creates the field. |
| `Text() string` | Current content. |
| `SetText(s string)` | Replaces the content. |
| `Submit(fn func(text string)) *Input` | Handler for Enter. |
| `Styles(st Styles) *Input` | Field styles. |
| `Focus()` | Requests keyboard focus. |

## Styles

| Property | Role |
| --- | --- |
| `Color` | Text color. |
| `Background` | Field fill. |
| `BorderWidth` / `BorderColor` / `BorderRadius` | Field frame; defaults to a 1px border. |
| `Padding` | Inner padding; defaults to 12x8. |

The field expands to the available width by default; set `Width` to fix it.
