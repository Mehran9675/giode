# Writing your own components

Anything with a `Layout` method is an element:

```go
type Element interface {
	Layout(gtx layout.Context) layout.Dimensions
}
```

So a type of your own can go anywhere a built-in can: in a `Box`, returned from a view,
inside another component. The built-in components are written exactly this way, from the
same pieces described here.

## Which approach

| You need | Use |
| --- | --- |
| A bit of raw Gio drawing or layout | `giode.Raw(func(gtx layout.Context) layout.Dimensions { ... })` |
| Reusable visuals, no state | A function returning elements, or a struct with a `Layout` method built from `Box`, `Text` and icons |
| Something clickable | [Pressable](components/pressable.md) around your visuals |
| Other interactive state | A struct holding Gio state (`widget.Editor`, `gesture.Drag`, ...), created once |

## A component step by step

1. **Make a struct for the state.** Users create it once with a `New...` function and lay it
   out every frame, like the built-ins.
2. **Build the visuals from elements** (`Box`, `Row`, `Text`, icons), dropping to raw Gio
   where needed.
3. **Accept `Styles`** as an optional last constructor argument, and apply your defaults
   with `styles.Merge(defaults, userStyles)` so users can override them.
4. **Return the struct from setters** (`OnChange`, `Label`, ...) so calls chain.
5. **Report changes** that happen without input with `kit.RequestFrame` (see
   [Rendering performance](performance.md#components-inside-a-memo)).

## Example: a star rating

```go
package rating

import (
	"gioui.org/layout"

	"github.com/mehran9675/giode/components/icon"
	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
)

var star = icon.Path{Data: "M12 2L15 9L22 9.5L16.5 14L18 21L12 17L6 21L7.5 14L2 9.5L9 9Z"}

// Rating is five clickable stars.
type Rating struct {
	value    int
	stars    [5]*elements.Pressable
	onChange func(int)
	st       styles.Styles
}

func New(st ...styles.Styles) *Rating {
	r := &Rating{value: 3}
	if len(st) > 0 {
		r.st = st[0]
	}
	for i := range r.stars {
		n := i + 1
		r.stars[i] = elements.NewPressable(func(s elements.PressState) elements.Element {
			return r.star(n, s)
		}).Label("Rate " + string(rune('0'+n))).OnClick(func(elements.ClickEvent) {
			r.value = n
			if r.onChange != nil {
				r.onChange(n)
			}
		})
	}
	return r
}

func (r *Rating) Value() int { return r.value }

func (r *Rating) OnChange(fn func(int)) *Rating {
	r.onChange = fn
	return r
}

func (r *Rating) star(n int, s elements.PressState) elements.Element {
	col := r.st.Color
	if col == "" {
		col = "#facc15"
	}
	if n > r.value {
		col = kit.Alpha(col, 0.25) // unselected stars are dim
	}
	if s.Hovered {
		col = kit.Hovered(col)
	}
	return star.Icon(styles.Styles{Color: col, Width: 24})
}

// Layout makes *Rating an element.
func (r *Rating) Layout(gtx layout.Context) layout.Dimensions {
	st := styles.Merge(styles.Styles{Gap: 4}, r.st)
	return elements.Row(st, r.stars[0], r.stars[1], r.stars[2], r.stars[3], r.stars[4]).Layout(gtx)
}
```

```go
stars := rating.New().OnChange(func(v int) { log.Printf("rated %d", v) })

app.Run(func() giode.Element {
	return giode.Box(giode.Styles{Padding: giode.UniformInset(24)}, stars)
})
```

## Kit helpers

`components/kit` has the helpers the built-in components use:

| Helper | Use |
| --- | --- |
| `Hovered(c)`, `Pressed(c)`, `Scale(c, f)`, `Alpha(c, a)` | Derive hover/pressed shades and transparency from a color. |
| `RequestFrame(gtx)`, `RequestFrameAt(gtx, t)`, `MarkLive()` | Ask for another frame; tell enclosing memos the component changed. |
| `HoverTracker` | Embed in an element to give it `OnHover`, like `Text` and `Divider`. |
| `BlockDrag(gtx, size)`, `BlockDragHere(gtx)` | Keep the component clickable inside a draggable frameless window or title bar. |
| `DragArea(gtx, size)` | Mark an area as a window drag handle, for a custom title bar. |
| `Scrollbar` | The scrollbar used by `Scroll` and `VirtualList`. |
| `DrawGlyph(gtx, g, size, color)` | The chevrons, check mark and window-button glyphs. |
| `WriteClipboard(gtx, text)` | Copy text from inside `Layout`. |
| `ClampPos(win, size, pos, margin)`, `Margin` | Keep a popup inside the window. |

From `elements`, for layout:

| Helper | Use |
| --- | --- |
| `Measure(gtx, el)` | Size of an element without drawing it or taking input. |
| `SizedPx(w, h, el)`, `OffsetPx(off, el)` | Place an element by a size or offset you measured in pixels. |
| `PushColor(c)`, `ResolveColor(c)` | Pass a text color down to children, and resolve one. |
| `RecordSemantic(gtx, node, size)` | Show the component to automation such as the [MCP server](mcp.md). |

## Popups: portals

Menus, dropdowns and tooltips have to draw above everything, however deep they sit in the
tree. `elements.RegisterPortal` queues content to be drawn after the rest of the frame, near
a window position:

```go
anchor := elements.AbsolutePosition().Add(image.Pt(0, dims.Size.Y+4)) // just below me
elements.RegisterPortal(gtx, anchor, func() elements.Element { return panel })
```

When the content is drawn it's kept inside the window, flipping above the anchor when
there's no room below.

`RegisterPopup` does the same for content that should close when the user clicks
elsewhere. It puts a transparent layer under the content that calls `OnDismiss` on a press
outside it:

```go
elements.RegisterPopup(gtx, anchor, s, func() elements.Element { return list },
	elements.PopupOptions{OnDismiss: s.Close})
```

By default the outside press doesn't reach the UI below, as with native menus. Set
`PassThrough: true` to let it through, for example for suggestions under a text field.

**Positions inside containers.** `AbsolutePosition()` is the window position of the element
being laid out, as far as is known at that moment. A `Row`, `Stack`, `Scroll` or `Tabs`
only decides where a child goes after laying it out, so anchors (and semantic nodes)
recorded inside a child are moved there afterwards. This happens automatically for
everything registered with `RegisterPortal`, `RegisterPopup` and `RecordSemantic`. Two
cases need a hand:

- A position you compute now and use elsewhere, such as a trigger whose menu is laid out
  in another part of the tree: store it with `elements.TrackPosition(&field, pos)` so it's
  corrected too. For an already-final window position, like the pointer position of a
  right-click, use `RegisterPortalAt`, which nothing moves.
- A container of your own that places a child with `op.Offset` after laying it out: take
  `mark := elements.MarkPositions()` before the child's `Layout` and call
  `elements.ShiftPositions(mark, offset)` once the offset is known. A container that knows
  the offset beforehand can call `elements.PushOffset(offset)` around the child instead.

## Sharing a component

Put it in a package of your own that imports Giode's `elements`, `styles` and `kit`.
Popups should use portals so they work wherever the component is placed. See
[Architecture](architecture.md) for how the built-ins are organized.
