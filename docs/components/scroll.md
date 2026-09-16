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
| `ContentHeight() int` | Total content height, as measured on the last frame. |
| `ViewportHeight() int` | Visible viewport height, as measured on the last frame. |
| `Progress() float32` | How far scrolled, from 0 (top) to 1 (bottom); 0 when there's nothing to scroll. |
| `AtBottom() bool` | Whether the scroll is at (or past) its maximum offset - also true when the content already fits the viewport. |
| `Visible(top, height int) bool` | Whether the vertical range `[top, top+height)` - in the content's own coordinate space, e.g. a running offset tracked while building a long list - is at least partially within the current viewport. |
| `OnReachBottom(fn func()) *Scroll` | Handler fired once when the scroll reaches its maximum offset (edge-triggered - fires again only after scrolling away and back). |

## Detecting scroll position

An intersection-observer-like pair for infinite scroll and lazy rendering:

```go
list := giode.Scroll(func() giode.Element {
	items := make([]giode.Element, 0, loaded)
	for i, top := 0, 0; i < loaded; i, top = i+1, top+itemHeight {
		if !list.Visible(top, itemHeight) {
			items = append(items, placeholder) // skip expensive content off-screen
			continue
		}
		items = append(items, renderItem(i))
	}
	return giode.Stack(giode.Styles{}, items...)
}, giode.Styles{})

list.OnReachBottom(func() {
	loaded += pageSize // fetch/reveal the next page
})
```

`Visible` and `AtBottom` reflect the scroll offset as of the start of this Layout call,
one frame behind the content closure's own not-yet-measured output - fine for both
patterns above: a placeholder swaps to real content one frame after scrolling it into
view, and "reached bottom" is a one-off trigger, not a per-frame read.

## Styles

| Property | Role |
| --- | --- |
| `Background` | Container fill. |
| `ScrollBar` | Scrollbar customization: `Width`, `TrackColor`, `ThumbColor`, `Radius`, `MinThumbLength`. |

Interactive children inside the content keep receiving events. For very long lists prefer
virtualized scrolling (a future addition).
