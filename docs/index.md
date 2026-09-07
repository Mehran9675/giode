# Giode

Giode is a UI library built on [Gio](https://gioui.org) that brings retained-mode ease to
immediate-mode UI. Interfaces are described with composable functions, styled with CSS-like
properties and wired up with plain Go callbacks — no webview, full performance, and raw Gio
stays reachable everywhere.

## Features

- **Window config** — size, min/max, resizable, frameless, draggable, fullscreen, centered,
  always-on-top, plus runtime actions (minimize, maximize, center, raise, close).
- **Elements** — `Box`, `Flex`, `Text`, `Divider`, `Raw`, `Responsive`.
- **Components** — Button, Input, Checkbox, Slider, Progress, Spinner, Tabs, Dropdown,
  Drawer, Dialog, Scroll, Title, Image, and a right-click context Menu with window-bounds
  clamping.
- **CSS-like styling** — 40+ properties, each defined in its own file, extensible by design.
- **Routing** — `NewRouter`, `Push`/`Navigate`/`Back`, `Link`.
- **Breakpoints** — `Responsive` switches views when the available width crosses a threshold.
- **Icons** — Material set, SVG path data, IconVG, plain images.
- **Low-level escape hatches** — `Raw`, `RunRaw`, `App.Window()`.

## Documentation

- [Getting started](getting-started.md)
- [Window configuration](window.md)
- [Elements](elements.md)
- [Style properties](styles.md)
- [Components](components/)
- [Routing](routing.md)
- [Architecture](architecture.md)
- [Low-level access](low-level.md)

## Packages

The root package `giode` re-exports the whole API for single-import ergonomics. The
underlying packages stay importable directly:

| Package | Contents |
| --- | --- |
| `styles/properties` | One file per CSS-like property |
| `styles` | The `Styles` aggregate, `Merge`, `Hex` |
| `elements` | Stateless building blocks |
| `components/*` | Stateful components, one directory each |
| `router` | Navigation |
| `fonts` | Font/shaper registry |
