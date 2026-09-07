# Progress

Stateless determinate bar. Build it every frame.

```go
giode.Progress(0.35, giode.Styles{Color: giode.HexColor("#3b82f6")})
```

| Argument | Description |
| --- | --- |
| `value` | Progress in [0, 1]. |
| `st` | Optional styles. |

## Styles

| Property | Role |
| --- | --- |
| `Color` | Fill. |
| `Background` | Track. |
| `Height` | Bar height; defaults to 8px. |

The bar expands to the available width.
