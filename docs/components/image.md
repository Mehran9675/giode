# Image

Shows an `image.Image`, scaled to fit its box.

```go
photo, _, _ := image.Decode(file) // import _ "image/jpeg" or "image/png"

giode.Image(photo, giode.Styles{
	Width: 120, Height: 120,
	Fit:          giode.FitCover,
	BorderRadius: 12,
})
```

| Style | Role |
| --- | --- |
| `Width`, `Height` | Size in dp; the image's own size when unset. |
| `Fit` | `FitContain` (default) shows all of it, `FitCover` fills the box and crops, `FitStretch` distorts to fill, `FitNone` keeps the original size. |
| `BorderRadius` | Rounded corners. |
| `Opacity` | Transparency. |

Images are uploaded to the GPU once and reused, as long as you pass the same `image.Image`
value. If you change an image's pixels in place, call `properties.ForgetImage(img)` so the
new pixels are uploaded.

`OnHover` changes the styles while the pointer is over it (for example dimming it with
`Opacity`); see [hover styles](../elements.md#hover-styles).
