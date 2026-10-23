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

## Hover

Every stateless element with something to paint supports `OnHover`: `Box` (and `Row`/
`Stack`), `Text`, `Divider`, `Image`, `Progress` and `Spinner`. Chain it to apply styles
while the pointer is over the element:

```go
card := giode.Box(base, children...).OnHover(giode.Styles{
	Background:   "#1e40af",
	BorderColor:  "#3b82f6",
	BorderWidth:  1,
	BoxShadow:    giode.BoxShadow{Y: 4, Blur: 12, Color: "#00000060"},
	Cursor:       giode.CursorPointer,
})

label := giode.Text("Docs", giode.Styles{Color: "#94a3b8"}).OnHover(giode.Styles{
	Color:          "#f8fafc",
	TextDecoration: giode.TextDecorationUnderline,
	Cursor:         giode.CursorPointer,
})
```

Hover styles merge on top of the base styles for both layout and painting, and can set
anything the element's own `Styles` support — colors, borders, shadows, padding,
cursor, and so on (see each element's own page for which properties it actually
paints). The hover state lives in the element, so **create it once and reuse it** like
a stateful component; rebuilding it inside the view function resets the hover each
frame. `Hovered()` reports the state observed on the last frame.

`Responsive` and `Raw` don't support `OnHover` - they're structural (they just delegate
to another element or to raw Gio), not something that paints its own styles to merge
into. The stateful components (`Button`, `Input.*`, `Tabs`, ...) already shade
themselves on hover internally instead of going through `OnHover`.

Component authors: `OnHover` is built on `kit.HoverTracker`, an embeddable building
block for adding the same capability to a custom stateless element - see
[Extending the library](extending.md).

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

## Table

Rows of cells with column widths aligned across all of them - see
[table.md](components/table.md):

```go
giode.Table(giode.Styles{},
	giode.TableHeader(giode.Styles{}, giode.TableCell(giode.Styles{}, giode.Text("Name"))),
	giode.TableRow(giode.Styles{}, giode.TableCell(giode.Styles{}, giode.Text("Ada"))),
)
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
	giode.Breakpoint{MaxWidth: 600, Name: "compact", View: func() giode.Element { return phoneView }},
	giode.Breakpoint{MaxWidth: 1024, Name: "medium", View: func() giode.Element { return tabletView }},
	giode.Breakpoint{View: func() giode.Element { return desktopView }}, // fallback
)
```

Breakpoints are evaluated in order; the first whose `MaxWidth` is at least the available
width matches, so list the narrowest first. A breakpoint without `MaxWidth` always matches.
Only the matching view is built each frame.

`MaxWidth` is in dp, like every length in `Styles`, so a layout switches at the same
physical size on a desktop and on a 3x phone. It's compared with the width the parent
container makes available, not the window's. The older `Width` field is in physical
pixels and is deprecated; it applies only when `MaxWidth` is 0.

Content inside the chosen view can ask which breakpoint matched, without a flag passed
through every page:

```go
if giode.ActiveBreakpoint() == "medium" { ... }
if giode.Compact() { ... } // the innermost Responsive matched a breakpoint named "compact"
```

Both answer for the innermost enclosing `Responsive` and return `""`/`false` outside one.
Like `AbsolutePosition`, they're only meaningful while the view is being laid out. A
`Memo` inside the view rebuilds when the active breakpoint changes.
