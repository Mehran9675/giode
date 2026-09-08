# Elements

Elements are the stateless building blocks. They are cheap to create inside the view
function every frame.

## Box

The container. Every `Box` lays its children out as a flex container along
`FlexDirection` (row by default), with `Align`, `Justify` and `Gap` distribution:

```go
giode.Box(giode.Styles{
	FlexDirection: giode.FlexDirectionRow, // or FlexDirectionColumn
	Align:         giode.AlignCenter,
	Justify:       giode.JustifySpaceBetween,
	Gap:           8,
}, child1, child2)
```

Box also paints the background stack (color, image, border, shadow) and acts as the
containing block for absolutely positioned children.

## Row and Stack

Shorthand containers with the flex direction preset:

```go
giode.Row(giode.Styles{Gap: 8, Align: giode.AlignCenter}, a, b, c)     // horizontal
giode.Stack(giode.Styles{Gap: 8}, a, b, c)                            // vertical
```

`Row` sets `FlexDirection` to row, `Stack` to column. Both behave exactly like a
plain `Box`.

### Wrapping

`Wrap` controls whether children continue onto a new line when they overflow the main
axis (CSS flex-wrap). It is **on by default for `Row`** (and any `Box` with
`FlexDirection: FlexDirectionRow`) and off by default for `Stack`/column containers;
either can be overridden explicitly:

```go
giode.Row(giode.Styles{Gap: 8}, tag1, tag2, tag3, /* ...many more... */)
// narrower than the tags need? they wrap onto additional lines automatically.

giode.Row(giode.Styles{Gap: 8, Wrap: giode.WrapNoWrap}, a, b, c) // force single line
giode.Stack(giode.Styles{Wrap: giode.WrapWrap}, a, b, c)         // wrap a column too
```

Each child is measured once, so wrapping is safe with stateful children (buttons,
inputs, ...). The one limitation this implies: `FlexGrow` has no effect on children of
a wrapping container — a line's leftover space can't be redistributed after a child
has already been laid out, so `Justify`/`Align` still space and align items, they just
don't grow them.

If a wrapped `Row` ends up taller than its container (vertical overflow), it does not
scroll on its own — wrap it in [`Scroll`](components/scroll.md), which is a stateful
component (create it once, like `Button`) and has a fully customizable scrollbar:

```go
list := giode.Scroll(func() giode.Element {
	return giode.Row(giode.Styles{Gap: 8}, tag1, tag2, tag3 /* ... */)
}) // once

giode.Box(giode.Styles{Height: 200}, list) // in the view
```

## Text

```go
giode.Text("Hello", giode.Styles{Color: giode.HexColor("#fff"), FontSize: 18})
```

The styles argument is optional. Supports text alignment, weight, style, family,
line-height, max-lines and decoration.

## Divider

A thin horizontal line. Color defaults to a subtle gray, height to 1px:

```go
giode.Divider(giode.Styles{Color: giode.HexColor("#334155")})
```

## Raw

Drops down to plain Gio at any point of the tree:

```go
giode.Raw(func(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{Size: image.Pt(10, 10)}
})
```

## Responsive

Simplified media queries: renders different views depending on the available width.

```go
giode.Responsive(
	giode.Breakpoint{Width: 600, View: func() giode.Element { return phoneView }},
	giode.Breakpoint{Width: 1024, View: func() giode.Element { return tabletView }},
	giode.Breakpoint{View: func() giode.Element { return desktopView }}, // fallback
)
```

Breakpoints are evaluated in order; the first with `width >= available` matches. A `Width`
of zero always matches. Only the matching view is built each frame. Width is measured in
the layout pixel space of the parent container.
