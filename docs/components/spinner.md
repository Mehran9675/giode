# Spinner

A spinning loading indicator.

```go
giode.Spinner(giode.Styles{Color: "theme.primary", Width: 24})
```

| Style | Role |
| --- | --- |
| `Color` | Arc color. |
| `Width`, `Height` | Size, 32dp by default. Set only `Width` for a square. |

It's stateless: build it in the view only while something is loading. While it's on screen
it keeps the window redrawing; when it's gone, redrawing stops.

`OnHover` restyles it while the pointer is over it; see [hover styles](../elements.md#hover-styles).
