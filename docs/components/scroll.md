# Scroll

Vertical scroll container. It implements `Element`, so it drops directly into the
tree like any other element — no raw Gio required. content is rebuilt every frame
and may be taller than the viewport. The mouse wheel, touch drags and the scrollbar
thumb all scroll.

```go
scroller := giode.Scroll(func() giode.Element { return content }) // once

// in the view:
giode.Box(giode.Styles{Height: 300}, scroller)

// fully customized scrollbar:
custom := giode.Scroll(func() giode.Element { return content }, giode.Styles{
	ScrollBar: giode.ScrollBar{
		Width: 12, Radius: 6, MinThumbLength: 40,
		TrackColor: "#1e293b", ThumbColor: "#3b82f6",
	},
})
```

## API

| Member | Description |
| --- | --- |
| `Scroll(content func() Element, st ...Styles) *Scroll` | Creates the container around content, rebuilt every frame. The styles argument is optional. |
| `ScrollTo(v float32)` | Sets the vertical offset in px. |
| `Offset() float32` | Current offset. |

## Styles

| Property | Role |
| --- | --- |
| `Background` | Container fill. |
| `ScrollBar` | `{Width, Radius, MinThumbLength int; TrackColor, ThumbColor color}` — every part of the bar. All fields are optional; unset ones fall back to sane defaults (8px wide, a translucent gray track, a lighter gray thumb, pill-rounded, 20px minimum thumb length). |

Interactive children inside the content keep receiving events. For very long lists prefer
virtualized scrolling (a future addition).
