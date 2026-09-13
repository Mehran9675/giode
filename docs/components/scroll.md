# Scroll

Vertical scroll container. The content closure is rebuilt every frame and may be taller
than the viewport. The mouse wheel, touch drags and the scrollbar thumb all scroll.

```go
scroller := giode.Scroll(func() giode.Element {
	// build the content fresh every frame
	return giode.Stack(giode.Styles{Gap: 8}, items...)
}, giode.Styles{ScrollBar: giode.ScrollBar{ThumbColor: "#3b82f6"}})

// in the view:
giode.Box(giode.Styles{Height: 300}, scroller)
```

## API

| Member | Description |
| --- | --- |
| `Scroll(content func() Element, st ...Styles) *Scroll` | Creates the container. |
| `ScrollTo(v float32)` | Sets the vertical offset in px. |
| `Offset() float32` | Current offset. |

## Styles

| Property | Role |
| --- | --- |
| `Background` | Container fill. |
| `ScrollBar` | Scrollbar customization: `Width`, `TrackColor`, `ThumbColor`, `Radius`, `MinThumbLength`. |

Interactive children inside the content keep receiving events. For very long lists prefer
virtualized scrolling (a future addition).
