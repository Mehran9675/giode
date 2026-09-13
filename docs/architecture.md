# Architecture

## Layout

```
giode.go / app.go / config.go / drag.go / reexport.go
    Root package: App, Config, Run, re-exports.
    singleinstance*.go  Single-instance lock (per-platform process check)
    icon*.go            Runtime window icon (Windows)
styles/
    styles.go        Styles aggregate
    merge.go         Merge(base, over)
    properties/      One file per CSS property (colors resolve via calcColor)
elements/
    Stateless building blocks (Box, Row, Stack, Text, Divider, Raw, Responsive)
components/
    One directory per component, one file per function or struct
    kit/               Public helpers for custom components
router/
    Navigation
fonts/
    Font and shaper registry
examples/
    hello, showcase, custom
```


## Adding a property

1. Create `styles/properties/<name>.go` with the type, constants and resolution logic.
2. Add the field to `Styles` in `styles/styles.go`.
3. Add the overlay rule to `Merge` in `styles/merge.go`.
4. Consume it in `elements` (or the relevant component).
5. Re-export constants from `reexport.go` if useful.

## Adding a component

1. Create `components/<name>/` with one file per function/struct.
2. Stateful components hold Gio widgets (`widget.Clickable`, `widget.Editor`, gestures)
   and are created once with a `New` constructor.
3. Build visuals out of `elements` (dogfooding) and raw Gio where needed.
4. Accept `styles.Styles` where sensible, using `styles.Merge` for defaults.
5. Use the `components/kit` helpers (color shading, overlay clamping) — the same
   helpers the built-in components use.
6. Re-export from `reexport.go`.

See [extending.md](extending.md) for the full recipe for custom components.

## State model

Immediate mode rebuilds the tree every frame, so persistent state lives in the stateful
handles, never in the view function. Stateless elements keep no state at all.
