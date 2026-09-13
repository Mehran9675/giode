# Elements

Elements are the stateless building blocks. They are cheap to create inside the view
function every frame.

## Box

The flex container. Children are laid along `FlexDirection` (row by default) with
`Align`, `Justify`, `Gap` and `Wrap`:

```go
giode.Box(giode.Styles{
	FlexDirection: giode.FlexDirectionColumn,
	Align:         giode.AlignCenter,
	Gap:           12,
	Padding:       giode.UniformInset(16),
}, child1, child2)
```

Box also paints the background stack (color, image, border, shadow) and acts as the
containing block for absolutely positioned children.

## Row and Stack

Shorthand containers with the flex direction preset:

```go
giode.Row(giode.Styles{Gap: 8, Align: giode.AlignCenter}, a, b, c) // horizontal
giode.Stack(giode.Styles{Gap: 8}, a, b, c)                        // vertical
```

`Row` sets `FlexDirection` to row, `Stack` to column.

## Wrap

By default a `Row` wraps children onto new lines when they overflow (`WrapWrap`), while
a `Stack` keeps them on a single line (`WrapNoWrap`). Override with the `Wrap` property:

```go
giode.Row(giode.Styles{Wrap: giode.WrapNoWrap, Gap: 8}, chips...)
giode.Stack(giode.Styles{Wrap: giode.WrapWrap, Gap: 8}, chips...)
```

Wrapped lines are stacked on the cross axis and separated by `Gap`.

## Text

```go
giode.Text("Hello", giode.Styles{Color: "#fff", FontSize: 18})
```

The styles argument is optional. Supports text alignment, weight, style, family,
line-height, max-lines and decoration.

## Divider

A thin horizontal line. Color defaults to a subtle gray, height to 1px:

```go
giode.Divider(giode.Styles{Color: "#334155"})
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
