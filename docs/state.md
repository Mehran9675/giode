# Global state

`giode.State[T]` is a minimal, generic observable value: read and write it from
anywhere, including a goroutine outside the render loop (a network response, a timer),
not just from within a `Layout` call.

```go
counter := giode.NewState(0)

btn := giode.Button("Increment").OnClick(func() {
	counter.Update(func(v int) int { return v + 1 })
})

// elsewhere - a goroutine, a network callback, a timer:
go func() {
	result := fetchSomething()
	counter.Set(result)
}()

// in the view:
giode.Text(fmt.Sprintf("Count: %d", counter.Get()))
```

Unlike every other piece of state in giode (a `Button`'s click count, a `Scroll`'s
offset), which lives inside one component and only ever changes from within a `Layout`
call, `State` is meant to be shared across the app and written from anywhere - `Set`/
`Update` request a redraw themselves (via the app window's `Invalidate`, wired up
automatically by `App.New`), so a background goroutine's update shows up without
needing its own `gtx.Execute(op.InvalidateCmd{})`.

## API

| Member | Description |
| --- | --- |
| `NewState[T](initial T) *State[T]` | Creates a State holding `initial`. |
| `Get() T` | Current value. |
| `Set(v T)` | Replaces the value, notifies watchers, requests a redraw. |
| `Update(fn func(current T) T)` | Atomically replaces the value with `fn(current)` - safe for a read-modify-write under concurrent access (e.g. an increment) where a separate `Get` then `Set` could race. |
| `Watch(fn func(T)) (unsubscribe func())` | Registers a callback for every future `Set`/`Update`. |

`State` is safe for concurrent use - `Get`, `Set`, `Update` and `Watch` can all be called
from multiple goroutines at once.

## When to reach for it

Most state still belongs inside one stateful component (a counter local to one page,
say) - a `State` value's job is specifically being shared across the view, or written
from outside the render loop:

- A value several unrelated parts of the view need to read (a signed-in user, a color
  theme flag, a shopping cart count in a header badge).
- A value updated by a background goroutine - the case its own redraw-triggering exists
  for.

For UI-only interactions confined to one component (a click count next to one button, a
dropdown's own selection), the component's own field is simpler - see
[Extending the library](extending.md).
