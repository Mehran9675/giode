# Architecture

## Layout

```
giode.go / app.go / config.go / drag.go / reexport.go
    Root package: App, Config, Run, re-exports.
styles/
    styles.go        Styles aggregate
    merge.go         Merge(base, over)
    hex.go           HexColor
    properties/      One file per CSS property
elements/
    Stateless building blocks (Box, Flex, Text, Divider, Raw, Responsive)
components/
    One directory per component, one file per function or struct
    internal/        Shared helpers (overlay clamping, color math)
router/
    Navigation
fonts/
    Font and shaper registry
examples/
    hello, showcase
```

## The rule of one

- Every CSS property is defined in exactly one file under `styles/properties/`.
- Every component lives in its own directory; each function and struct in its own file.
- Complex properties or components may spread across several files in their directory.

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
5. Re-export from `reexport.go`.

## State model

Immediate mode rebuilds the tree every frame, so persistent state lives in the stateful
handles, never in the view function. Stateless elements keep no state at all.
