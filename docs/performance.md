# Rendering performance

The view function rebuilds the whole UI every frame. Gio only asks for a frame when
something happens (input, a resize, a state change, a running animation), so an idle app
does no work. For large UIs that change in small places, `Memo` avoids rebuilding the parts
that didn't change.

## Memo

`Memo` caches what a part of the UI drew. On later frames it replays the cached drawing
instead of rebuilding, which costs almost nothing however big that part is.

```go
giode.Memo("sidebar", func() giode.Element {
	return sidebar(user.Name)
}, user.Name) // dependencies: plain values the builder reads
```

The first argument is a key that identifies the cache from frame to frame. It must be
unique among the memos laid out in one frame; a memo that isn't laid out in a frame is
discarded. Inside a component of your own, keep a `*MemoCache` from `NewMemoCache()` and
call its `Layout` method instead of using a key.

A memo rebuilds when:

- one of the listed dependencies changed (compared with `==`, or by deep comparison for
  slices and maps);
- a `State` read while building changed (tracked automatically);
- the available size or the screen density changed;
- the theme changed (`SetMode`, `SetPalette`);
- the inherited text color changed;
- the pointer is over it, so hover, presses and scrolling inside are handled normally;
- a component inside is animating, focused, being dragged or just handled an event;
- `InvalidateCaches()` was called.

Values the builder reads from ordinary variables must be listed as dependencies, or the memo
won't notice when they change. Memos can be nested: a list memo can hold one memo per row,
so changing one row rebuilds only that row.

## Components inside a memo

The built-in components handle this themselves. A component of your own that changes
without input, such as an animation or a countdown, has to report it, or the memo keeps
replaying the old drawing:

```go
kit.RequestFrame(gtx)      // animating: draw again soon
kit.RequestFrameAt(gtx, t) // something changes at time t
kit.MarkLive()             // changed for another reason
```

Use `kit.RequestFrame` instead of Gio's `gtx.Execute(op.InvalidateCmd{})`, which memos can't
see. When you drive a component from code (for example `button.Click()` from a timer), call
`InvalidateCaches()` so cached parts pick up the change; `app.RunOnUI` and `app.Post` do
this for you.

## Other built-in savings

Without any memo, Giode still avoids repeated work:

- colors are parsed once and cached;
- SVG files and path data are parsed once; IconVG data is decoded once;
- images are uploaded to the GPU once instead of every frame;
- `VirtualList` and `TreeView` build only the rows on screen;
- `State.SetIfChanged` skips the redraw when the value didn't change.
