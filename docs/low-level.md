# Low-level access

Raw Gio stays reachable at every layer.

## Raw elements

Drop to plain Gio at any point in the tree:

```go
giode.Box(giode.Styles{},
	giode.Raw(func(gtx layout.Context) layout.Dimensions {
		// any Gio layout, widget, paint or input code
		return layout.Dimensions{}
	}),
)
```

## Raw run loop

Skip the element model entirely:

```go
app := giode.New(giode.Config{...})
app.RunRaw(func(gtx layout.Context) layout.Dimensions {
	// a complete Gio frame
	return layout.Dimensions{}
})
```

## The window

`app.Window()` returns the underlying `*app.Window` for options, actions and events not
covered by `Config` and the action methods.

## GLFW and embedded windows

Only the root `App` layer depends on `gioui.org/app`. Every element, component, style
property, the router, icons and fonts are window-agnostic: they only need a
`layout.Context`. Rendering through a custom window (GLFW, Ebitengine, an embedded
surface) therefore works unchanged — build the context yourself each frame:

```go
// per frame, mirroring the gioui.org/app loop:
a.ops.Reset()
gtx := layout.Context{
	Ops:         &a.ops,
	Now:         time.Now(),
	Source:      a.router.Source(),          // input.Router fed by your window callbacks
	Metric:      unit.Metric{PxPerDp: scale, PxPerSp: scale},
	Constraints: layout.Exact(image.Pt(fbw, fbh)), // framebuffer size
}
view().Layout(gtx)                          // any Giode view works as-is
_ = gpuCtx.Frame(gtx.Ops, gpu.OpenGLRenderTarget{}, image.Pt(fbw, fbh))
a.router.Frame(gtx.Ops)                     // hit-testing for the next frame
window.SwapBuffers()
```

Notes for custom loops:

- Queue pointer/key events into `input.Router` in the same pixel space as the layout
  constraints; priority, pass-through and grab semantics are all handled by the router.
- The context menu works manually: call `menu.Context(gtx)` **before** the view and
  `menu.Layout(gtx)` **after** it (the equivalent of `SetContextMenu`).
- `Config.Draggable` relies on Gio's window driver (`system.ActionInputOp(ActionMove)`)
  and is a no-op under GLFW — use native window dragging there instead.
- The `Config` window options and the `App` action methods (fullscreen, minimize,
  center, always-on-top, …) are `app.Window`-only.
- Animations (Spinner, Drawer, Dialog) request frames with `op.InvalidateCmd`. A
  continuous render loop ignores that naturally; an event-driven loop must poll
  `router.WakeupTime()` or the animations freeze.

## Direct package imports

The root package is a re-export convenience. Everything is also importable from its
natural package:

```go
import (
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/components/menu"
	"github.com/mehran9675/giode/styles/properties"
)
```
