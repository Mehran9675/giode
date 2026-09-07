# Scroll

Vertical scroll container. The content is rebuilt every frame and may be taller than the
viewport. The mouse wheel, touch drags and the scrollbar thumb all scroll.

```go
scroller := giode.NewScroll()

// in the view:
giode.Box(giode.Styles{Height: 300}, giode.Raw(func(gtx layout.Context) layout.Dimensions {
	return scroller.Layout(gtx, content)
}))
```

## API

| Member | Description |
| --- | --- |
| `NewScroll() *Scroll` | Creates the container. |
| `Styles(st Styles) *Scroll` | Styles. |
| `ScrollTo(v float32)` | Sets the vertical offset in px. |
| `Offset() float32` | Current offset. |
| `Layout(gtx, content Element)` | Lays out the clipped, scrolled content. |

## Styles

| Property | Role |
| --- | --- |
| `Background` | Container fill. |
| `BorderColor` | Scrollbar thumb. |

Interactive children inside the content keep receiving events. For very long lists prefer
virtualized scrolling (a future addition).
