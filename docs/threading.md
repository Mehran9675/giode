# Threading

Everything that touches the UI runs on one goroutine, the **UI goroutine**:

- the view function passed to `App.Run`;
- every `Layout` method;
- content functions (`Scroll`, `VirtualList`, `TreeView`, `Tabs`, ...);
- event handlers (`OnClick`, `OnChange`, `Submit`, menu actions, ...).

So the rule is simple: if your data is only changed on the UI goroutine, it needs no locks.

## Results from other goroutines: Post

Work done in the background (a network request, reading a file, a timer) hands its result
back with `app.Post`:

```go
go func() {
	items, err := fetchItems()
	app.Post(func() {
		model.items, model.err = items, err // runs on the UI goroutine
	})
}()
```

`Post` returns immediately. The function runs before the next frame, which `Post` requests,
and functions run in the order they were posted. `app.RunOnUI(ctx, fn)` does the same but
waits until `fn` has run (or `ctx` ends).

`State` values are safe to set from any goroutine directly, and redraw when set.

## Changes during a frame: AfterLayout

Handlers run while the frame is being built. If a handler changes data that the rest of
the frame still reads, such as deleting the list item whose button was just clicked, defer
the change until the frame is built:

```go
remove.OnClick(func() {
	giode.AfterLayout(func() { model.delete(id) })
})
```

`AfterLayout` functions run on the UI goroutine once the frame is built, and another frame
follows right away so the change shows.

## Sizes without a layout context

Styles use dp. Code without a layout context can convert using the last frame's display
scale: `giode.DisplayScale()`, `giode.DpToPx(v)` and `giode.PxToDp(px)`.
