# State

Most state lives in components: a text field holds its text, a `Scroll` its offset. For data
that several parts of the UI read, or that changes outside the UI (a download finishing, a
timer firing), use `State`.

## State[T]

`State[T]` holds a value of any type. Reading and writing are safe from any goroutine, and
every write redraws the window.

```go
var counter = giode.NewState(0)

inc := giode.Button("Add one").OnClick(func() {
	counter.Update(func(v int) int { return v + 1 })
})

go func() {
	for range time.Tick(time.Second) {
		counter.Update(func(v int) int { return v + 1 })
	}
}()

// in the view:
giode.Text(fmt.Sprintf("Count: %d", counter.Get()))
```

| Method | Description |
| --- | --- |
| `NewState[T](initial T) *State[T]` | Creates a state. |
| `Get() T` | Current value. |
| `Set(v T)` | Replaces the value, notifies watchers, redraws. |
| `SetIfChanged(v T) bool` | `Set`, skipped when `v` equals the current value. |
| `Update(fn func(T) T)` | Read, modify and write in one atomic step. Use it for increments and appends instead of `Get` then `Set`, which can race. |
| `Watch(fn func(T)) func()` | Calls `fn` after every change, in registration order. Returns a function that stops watching. |
| `Version() uint64` | Increases with every change. |

When a slice or map is stored, `Update` should return a new one rather than modifying the
old one in place, so code still holding the old value doesn't see it change.

**With Memo.** A [`Memo`](performance.md) records which states it read while building, and
rebuilds when one of them changes. No extra dependency list is needed for state.

## Named state

`NamedState` registers a state under a name so tools can find it: the
[MCP server](mcp.md)'s `state_list` and `state_set` read and change it. The value must be
JSON-encodable for that.

```go
var todos = giode.NamedState("todos", []string{})
```

## Per-item state for dynamic lists

Components have to survive from one frame to the next, which is easy when you create them
once. It's harder for a list whose items come and go at runtime. `Keyed` keeps one
component per item, looked up by a stable key:

```go
func fileRows(files []File) []giode.Element {
	rows := make([]giode.Element, 0, len(files))
	for _, f := range files {
		row := giode.Keyed("files", f.ID, func() *giode.PressableElement {
			return giode.Pressable(func(s giode.PressState) giode.Element {
				return fileRow(f, s)
			}).OnClick(func(giode.ClickEvent) { open(f) })
		})
		rows = append(rows, row)
	}
	return rows
}
```

- Key by ID, not by position. With positions as keys, deleting an item hands its state to
  the item that moves into its place.
- The first argument (`"files"`) is a scope that keeps different lists apart when their
  keys could overlap.
- An entry that hasn't been asked for in 60 frames is dropped, so removed items don't
  leak. Frames only happen when something changes, so that isn't a fixed time.
  `SetKeyedRetention(n)` changes the limit and `ForgetKeyed(scope, key)` drops an entry
  right away.

Inside a component of your own, a `KeyedCache` does the same without a global scope:

```go
type fileList struct {
	rows *giode.KeyedCache[string, *giode.PressableElement]
}

l.rows = giode.NewKeyedCache[string, *giode.PressableElement]().
	Retain(120).
	OnEvict(func(id string, p *giode.PressableElement) { /* release resources */ })

row := l.rows.Get(f.ID, func() *giode.PressableElement { ... })
```

`KeyedCache` also has `Peek`, `Delete`, `Len` and `Range`.
