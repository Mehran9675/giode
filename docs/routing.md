# Routing

The router switches between views by path. It's an element, so place it wherever the
current page should appear.

```go
router := giode.NewRouter(
	giode.Route{Path: "/", View: homeView},
	giode.Route{Path: "/settings", View: settingsView},
).NotFound(notFoundView)

nav := giode.Row(giode.Styles{Gap: 12},
	router.Link("/", "Home"),
	router.Link("/settings", "Settings"),
)

app.Run(func() giode.Element {
	return giode.Stack(giode.Styles{Width: -1, Height: -1}, nav, router)
})
```

A view is a `func() giode.Element`. Only the current route's view is built.

| Method | Description |
| --- | --- |
| `NewRouter(routes ...Route) *Router` | Creates the router. The first route is shown at start. |
| `Navigate(path)` | Shows `path`, replacing the current history entry. |
| `Push(path)` | Shows `path`, keeping the current one in the history. |
| `Back() bool` | Returns to the previous path. Reports false if there is none. |
| `Current() string` | The current path. |
| `NotFound(view) *Router` | View for paths without a route. |
| `Link(path, text, st...) Element` | Text that navigates to `path` when clicked. |

Call `Push` or `Navigate` from any event handler:

```go
save := giode.Button("Save").OnClick(func() {
	store.Save()
	router.Push("/")
})
```

Components used by a view keep their state only if they're created outside it. Create them
once, or with [`Keyed`](state.md#per-item-state-for-dynamic-lists).
