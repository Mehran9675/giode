# Elements

Elements are the stateless building blocks. They are cheap to create inside the view
function every frame.

## Box

The container. `Display` selects the layout mode:

```go
// Block: children stack vertically, filling the width.
giode.Box(giode.Styles{Padding: giode.UniformInset(12)}, child1, child2)

// Flex: children laid along an axis.
giode.Box(giode.Styles{
	Display:   giode.DisplayFlex,
	Direction: giode.DirRow,          // or DirColumn
	Align:     giode.AlignCenter,
	Justify:   giode.JustifySpaceBetween,
	Gap:       8,
}, child1, child2)
```

Box also paints the background stack (color, image, border, shadow) and acts as the
containing block for absolutely positioned children.

## Flex

Shorthand for `Box` with `Display` forced to flex.

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
