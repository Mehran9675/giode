# Elements

Elements are the stateless building blocks. Build them inside the view function, every
frame; they cost very little. Anything with a `Layout(gtx layout.Context) layout.Dimensions`
method is an element, so components and your own types fit anywhere an element does.

## Box, Row and Stack

`Box` is a flex container: it lays its children along `FlexDirection` (a row by default),
spaced by `Gap`, aligned by `Align` and `Justify`.

```go
giode.Box(giode.Styles{
	FlexDirection: giode.FlexDirectionColumn,
	Align:         giode.AlignCenter,
	Gap:           12,
	Padding:       giode.UniformInset(16),
	Background:    "theme.surface",
	BorderRadius:  8,
}, title, body)
```

`Row` and `Stack` are a `Box` with the direction preset:

```go
giode.Row(giode.Styles{Gap: 8, Align: giode.AlignCenter}, icon, label) // horizontal
giode.Stack(giode.Styles{Gap: 8}, header, list, footer)               // vertical
```

A box also paints its own background, gradient, border and shadow, and positions
absolutely positioned children (see [Style properties](styles.md#absolute-positioning)).

**Sizing.** A box is as big as its content unless told otherwise. `Width: -1` or
`Height: -1` fills the available space, a positive value fixes the size in dp, and
`FlexGrow` shares the space left over in the parent's direction:

```go
giode.Row(giode.Styles{Width: -1, Wrap: giode.WrapNoWrap},
	giode.Box(giode.Styles{FlexGrow: 1}, giode.Text("Title")), // takes the spare width
	closeButton,                                              // pushed to the right edge
)
```

**Wrapping.** A `Row` wraps its children onto new lines when they don't fit; a `Stack`
doesn't. Set `Wrap: giode.WrapNoWrap` or `giode.WrapWrap` to change that. In a wrapping
row, `FlexGrow` has no effect, so use `WrapNoWrap` for toolbars that should stretch.

## Text

```go
giode.Text("Hello", giode.Styles{Color: "theme.text", FontSize: 18, FontWeight: giode.FontWeightBold})
```

Text wraps to the available width. `MaxLines` limits the line count and
`TextOverflow: giode.TextOverflowEllipsis` ends cut-off text with "…". Text without a
`Color` takes the color of the nearest enclosing box that sets one, or black.

## Divider

A thin horizontal line across the available width: 1dp high and a muted gray by default.

```go
giode.Divider(giode.Styles{Color: "theme.border"})
```

## Hover styles

`Box`, `Row`, `Stack`, `Text`, `Divider`, `Image`, `Progress` and `Spinner` accept extra
styles applied while the pointer is over them:

```go
card := giode.Box(base, children...).OnHover(giode.Styles{
	Background: "theme.border",
	Cursor:     giode.CursorPointer,
})
```

The hover styles are merged over the base styles. Add `Transition` to the base styles to
animate the change. The element keeps the hover state, so create it once and reuse it, or
look it up with [`Keyed`](state.md#per-item-state-for-dynamic-lists) inside lists.
`Hovered()` reports the state.

## Clickable content

To make any content clickable, use [Pressable](components/pressable.md), or call
`OnClick` on a box:

```go
row := giode.Box(giode.Styles{Padding: giode.UniformInset(8)}, label).
	OnClick(func(e giode.ClickEvent) { selectRow(id) })
```

## Responsive

Shows different views depending on the available width:

```go
giode.Responsive(
	giode.Breakpoint{Width: 600, View: phoneView},   // up to 600dp
	giode.Breakpoint{Width: 1024, View: tabletView}, // up to 1024dp
	giode.Breakpoint{View: desktopView},             // anything wider
)
```

The first breakpoint whose `Width` (in dp) is at least the available width wins; a zero
`Width` always matches, so put it last. Only the chosen view is built.

## Memo

`Memo` skips rebuilding a part of the UI when nothing it depends on has changed. See
[Rendering performance](performance.md).

```go
giode.Memo("sidebar", func() giode.Element { return sidebar(user) }, user.ID)
```

## Raw

Raw Gio code anywhere in the tree:

```go
giode.Raw(func(gtx layout.Context) layout.Dimensions {
	// any Gio layout, paint or input code
	return layout.Dimensions{Size: image.Pt(10, 10)}
})
```

## Table

Rows of cells whose columns line up. See [Table](components/table.md).
