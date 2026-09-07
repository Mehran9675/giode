# Image

Stateless image with fit modes.

```go
giode.Image(img, giode.Styles{
	Width: 120, Height: 120,
	Fit: giode.FitCover,
	BorderRadius: 12,
})
```

| Argument | Description |
| --- | --- |
| `img` | Any `image.Image` (decode PNG/JPEG with the standard library). |
| `st` | Optional styles. |

## Styles

| Property | Role |
| --- | --- |
| `Width` / `Height` | Box size; natural size when unset. |
| `Fit` | `FitContain` (default), `FitCover`, `FitStretch`, `FitNone`. |
| `BorderRadius` | Rounded corners. |
