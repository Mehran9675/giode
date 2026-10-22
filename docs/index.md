# Giode

Giode is a desktop UI library for Go, built on [Gio](https://gioui.org). You describe the
interface with plain functions, style it with CSS-like properties and handle events with Go
callbacks. It renders natively through Gio (no web view), and raw Gio is available at every
level when you need it.

```go
app := giode.New(giode.Config{Title: "Hello", Resizable: true, Background: "theme.background"})
count := 0
inc := giode.Button("Add one").OnClick(func() { count++ })

app.Run(func() giode.Element {
	return giode.Stack(giode.Styles{Gap: 12, Padding: giode.UniformInset(24)},
		giode.Text(fmt.Sprintf("Count: %d", count)),
		inc,
	)
})
```

## What's included

- **Layout and styling.** `Box`, `Row` and `Stack` are flex containers. `Styles` covers
  sizes, padding, gaps, borders, rounded corners, shadows, gradients, transforms, text
  styling and hover transitions. Sizes are in dp, so layouts keep their physical size on
  HiDPI screens.
- **Components.** Buttons, text fields, checkboxes, sliders, selects, autocomplete, tabs,
  dialogs, drawers, menus, tooltips, accordions, pagination, tables, a virtualized list, a
  tree view, toasts, a custom window title bar, transitions and an app-shell layout.
- **Theming.** Light and dark palettes with semantic colors (`"theme.primary"`,
  `"theme.surface"`, ...) that follow `SetMode`.
- **State.** `State[T]` is an observable value you can set from any goroutine. `Keyed` gives
  each item of a changing list its own persistent state.
- **Efficient redraws.** `Memo` caches parts of the UI and rebuilds them only when what they
  depend on changes.
- **Window control.** Size limits, frameless windows, window state, single-instance lock,
  clipboard, file drop (Windows) and lifecycle callbacks.
- **MCP server.** One config option lets AI agents inspect and drive the running app.
- **Icons you choose.** No bundled icon set: load SVG files or path data from any set.

## Where to start

- [Getting started](getting-started.md): install and a first app.
- [Elements](elements.md) and [Style properties](styles.md): building and styling the UI.
- [Components](components/README.md): the stateful widgets.
- [State](state.md) and [Threading](threading.md): where data lives and who may change it.

## Packages

The root package `giode` re-exports everything, so one import is enough. The packages below
can also be imported directly.

| Package | Contents |
| --- | --- |
| `elements` | Box, Text, Pressable, Memo, Keyed and the other building blocks |
| `styles`, `styles/properties` | The `Styles` struct and one file per property |
| `components/*` | One package per component |
| `components/kit` | Helpers for writing your own components |
| `state` | `State[T]` |
| `router` | Navigation between views |
| `mcp` | The MCP server, usable on its own |
| `fonts` | Font registration |
