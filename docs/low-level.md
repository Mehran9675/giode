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
