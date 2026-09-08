# Extending the library

Everything in Giode is an `Element`:

```go
type Element interface {
	Layout(gtx layout.Context) layout.Dimensions
}
```

Any type with a `Layout` method is a first-class Giode component — it can be returned
from a view, placed inside a `Box`/`Row`/`Stack` and re-exported next to the built-ins. On top of
that interface, Giode offers a small kit of helpers that the built-in components
themselves are made of.

## Quick starts

| Need | Use |
| --- | --- |
| One-off raw Gio code | `giode.Raw(func(gtx layout.Context) layout.Dimensions { ... })` |
| Stateless visual | A struct with a `Layout` method built from `Box`/`Row`/`Stack`/`Text`/icons |
| Stateful interactive component | A handle type holding Gio state (`widget.Clickable`, `widget.Editor`, `gesture.*`), created once with a `New` constructor |

## The component recipe

1. **Define a handle type.** Hold all persistent state in it. Create it once before
   `Run` and lay it out every frame — never construct it inside the view function.

2. **Hold Gio state directly.** Gio's `widget.Clickable`, `widget.Bool`, `widget.Editor`
   and `gesture.*` types are the building blocks of interactivity. They carry their own
   event queues; just call their `Clicked`/`Update`/`Layout` methods each frame.

3. **Build the visuals out of elements** (`Box`, `Row`, `Stack`, `Text`, icons) and drop to raw
   Gio ops where needed. This is exactly how the built-in components are written.

4. **Take `styles.Styles` as a constructor argument** (optionally variadic, as `Text`
   and the built-in components do), not a chained setter, and merge defaults with
   `styles.Merge(styles.Styles{...defaults...}, userStyles)` so users can override the
   relevant properties.

5. **Chain your setters**: return the component itself from `OnClick`-style methods.

## Kit helpers

The `components/kit` package (re-exported from the root) provides the shared pieces
component authors need:

| Helper | Description |
| --- | --- |
| `kit.Scale(c, f)` | Scales a color's channels by `f`, clamped. |
| `kit.Hovered(c)` | Hover shading: slightly lighter. |
| `kit.Pressed(c)` | Pressed shading: slightly darker. |
| `kit.Alpha(c, a)` | Multiplies a color's alpha. |
| `kit.ClampPos(win, size, pos, margin)` | Keeps an overlay panel inside a window; flips above the anchor when it would overflow the bottom edge. |
| `kit.Margin` | Default overlay-to-edge distance (8px). |

## Worked example: a star rating

```go
package main

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/icon"
	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
)

var gold = "#fac815"

// Rating is a custom stateful component: five clickable stars.
type Rating struct {
	value    int
	stars    [5]widget.Clickable
	onChange func(v int)
	st       styles.Styles
}

func NewRating() *Rating {
	return &Rating{value: 3}
}

func (r *Rating) Value() int { return r.value }

func (r *Rating) OnChange(fn func(v int)) *Rating {
	r.onChange = fn
	return r
}

func (r *Rating) Styles(st styles.Styles) *Rating {
	r.st = st
	return r
}

// Layout implements elements.Element, so a *Rating can be placed
// anywhere an element is expected.
func (r *Rating) Layout(gtx layout.Context) layout.Dimensions {
	// Defaults first, user styles on top.
	st := styles.Merge(styles.Styles{Gap: 4, Color: gold}, r.st)

	var kids []elements.Element
	for i := range r.stars {
		kids = append(kids, &starEl{rating: r, index: i, st: st})
	}
	return elements.Row(st, kids...).Layout(gtx)
}

// starEl is a stateless element for one star; the rating owns the
// clickable state.
type starEl struct {
	rating *Rating
	index  int
	st     styles.Styles
}

func (s *starEl) Layout(gtx layout.Context) layout.Dimensions {
	click := &s.rating.stars[s.index]
	if click.Clicked(gtx) {
		s.rating.value = s.index + 1
		if s.rating.onChange != nil {
			s.rating.onChange(s.rating.value)
		}
	}
	return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		col := s.st.Color
		if s.index >= s.rating.value {
			col = kit.Alpha(col, 0.25) // unselected stars are dim
		}
		return icon.Material("star", styles.Styles{Color: col, Width: 24, Height: 24}).Layout(gtx)
	})
}
```

Usage:

```go
rating := NewRating().OnChange(func(v int) { log.Printf("rated %d", v) })

app.Run(func() giode.Element {
	return giode.Box(giode.Styles{Justify: giode.JustifyCenter, Width: -1, Height: -1}, rating)
})
```

The same recipe applies to any component: an editor-based field holds a
`widget.Editor`, a dropdown holds a `widget.Clickable` plus `kit.ClampPos` for its panel,
a toggle holds a `widget.Bool`.

## Publishing your component

- Keep the handle type and its constructor in a package of your own; import Giode's
  `elements`, `styles` and `kit`.
- To mimic the built-ins, re-export it from the root `giode` package of a fork or wrap
  your own package in a small facade.
- Overlays (menus, dropdowns, popovers) must be laid out last in the view to render above
  the rest of the UI.

See [architecture.md](architecture.md) for how the built-in components are organized.
