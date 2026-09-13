# Spinner

Stateless animated loading indicator.

```go
giode.Spinner(giode.Styles{Width: 24, Height: 24, Color: "#3b82f6"})
```

## Styles

| Property | Role |
| --- | --- |
| `Color` | Arc color. |
| `Width` / `Height` | Size; defaults to 32px. |

The spinner invalidates its own frames while visible; no state is required.
