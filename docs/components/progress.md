# Progress

A horizontal progress bar. It's stateless: build it in the view with the current value.

```go
giode.Progress(download.Fraction(), giode.Styles{Color: "theme.primary"})
```

`value` goes from 0 to 1. The bar fills the available width.

| Style | Role |
| --- | --- |
| `Color` | Filled part. |
| `Background` | Track. |
| `Height` | Thickness, 8dp by default. |

`OnHover` restyles the bar while the pointer is over it; see [hover styles](../elements.md#hover-styles).
