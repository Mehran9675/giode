# Rendering performance

Giode is immediate mode: the view function describes the whole UI every frame. That keeps
the API simple, and Gio only produces a frame when something happens — input, a resize, a
`State` change, an animation. Within a frame, Giode avoids redoing work for the parts of the
UI that didn't change.

## Memo: only rebuild what changed

Wrap a part of the UI in `Memo` and it's rebuilt only when something it depends on changed;
every other frame replays the drawing recorded last time, at almost no cost regardless of
its size:

```go
giode.Memo("sidebar", func() giode.Element {
	return sidebar(user.Name)
}, user.Name) // deps: plain values the subtree reads
```

A memo rebuilds when:

- one of its **deps** changed (compared with `==`, or `reflect.DeepEqual` for slices, maps
  and other uncomparable values);
- a **State** read inside it changed — tracked automatically, no deps needed;
- its **constraints** (available size) or the display density changed;
- the **theme** changed (`SetMode`, `SetPalette`);
- a widget inside is **live**: animating (a spinner, a hover transition), focused (a text
  field), being dragged, or just handled an event;
- the **pointer** is over it, so hover, presses and scrolling are always handled live;
- `InvalidateCaches()` was called.

The key identifies the cache across frames and must be unique among the memos laid out in
a frame. A memo not laid out in a frame is released. Inside a stateful component, keep a
`*MemoCache` (from `NewMemoCache`) instead of using a key.

Memos nest: a list memo can contain one memo per row, so changing one row rebuilds only
that row.

### Rules for widgets inside a memo

Built-in components handle this themselves. A custom component inside a memo that changes
without input — an animation, a timer — must say so:

```go
kit.RequestFrame(gtx)       // animate: next frame ASAP, and keep enclosing memos live
kit.RequestFrameAt(gtx, t)  // timed change
kit.MarkLive()              // changed for another reason
```

Use `kit.RequestFrame` instead of `gtx.Execute(op.InvalidateCmd{})`, which enclosing memos
can't see. Values read from plain variables must be passed as deps. Driving a component
programmatically from outside a frame (`Button.Click()`) inside a cached subtree needs
`InvalidateCaches()`; `App.RunOnUI` does that for you.

## Built-in savings

Without any memo, Giode still avoids repeated work every frame:

- colors are parsed once and cached;
- SVG documents and path data are parsed once; IconVG is decoded once;
- images keep their GPU upload instead of being uploaded every frame;
- boxes without absolutely positioned or z-indexed children don't allocate to sort them;
- `State.SetIfChanged` skips redraws when a value didn't actually change.
