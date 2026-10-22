# Low-level access

Raw Gio is available at every level.

## Raw elements

```go
giode.Box(giode.Styles{},
	giode.Raw(func(gtx layout.Context) layout.Dimensions {
		// any Gio layout, widget, paint or input code
		return layout.Dimensions{}
	}),
)
```

## A raw view

`RunRaw` skips the element tree and gives you each frame's context directly:

```go
app.RunRaw(func(gtx layout.Context) layout.Dimensions {
	// a complete Gio frame
	return layout.Dimensions{}
})
```

## The Gio window

`app.Window()` returns Gio's `*app.Window` for options and actions Giode doesn't wrap.
`app.NativeHandle()` and `app.NativeView()` give the platform window; see
[Window configuration](window.md#native-access).

## Your own window or render loop

Only the root package depends on `gioui.org/app`. Elements, components, styles, the router,
icons and fonts only need a `layout.Context`, so they also render in a window you manage
yourself (GLFW, an embedded surface). Each frame, do what `App.Run` does:

```go
ops.Reset()
gtx := layout.Context{
	Ops:         &ops,
	Now:         time.Now(),
	Source:      router.Source(), // an input.Router fed by your window's events
	Metric:      unit.Metric{PxPerDp: scale, PxPerSp: scale},
	Constraints: layout.Exact(image.Pt(width, height)),
}

elements.BeginFrame()
elements.SetWindowSize(gtx.Constraints.Max)
elements.SetMetric(gtx.Metric)
view().Layout(gtx)
elements.DrainPortals(gtx) // popups: menus, tooltips, dropdowns, toasts
if elements.RunAfterLayout() {
	requestRedraw()
}
renderer.Frame(gtx.Ops) // e.g. gpu.GPU.Frame with your render target
router.Frame(gtx.Ops)   // hit testing for the next frame's input
elements.EndFrame()
```

Notes:

- Feed pointer and key events to the `input.Router` in pixels, in the same space as the
  constraints.
- `State` writes call the function registered with `state.SetInvalidator`; set it to your
  redraw request. Do the same for the theme with `properties.SetThemeChangeHook`.
- Animations request frames with Gio's invalidate command. In an event-driven loop, check
  `router.WakeupTime()` after each frame and schedule the next one, or animations stop.
- A context menu needs `menu.Context(gtx)` before the view and `menu.Layout(gtx)` after
  it; `App.SetContextMenu` does this for you.
- `Config.Draggable`, `TitleBar`'s drag area and the window actions rely on Gio's window
  and do nothing elsewhere.

## Importing packages directly

The root package only re-exports. Everything can be imported from its own package:

```go
import (
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/components/menu"
	"github.com/mehran9675/giode/styles/properties"
)
```
