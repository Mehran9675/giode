# Giode

Giode is a UI library built on [Gio](https://gioui.org) that brings retained-mode ease to
immediate-mode UI. Interfaces are described with composable functions, styled with CSS-like
properties and wired up with plain Go callbacks — no webview, full performance, and raw Gio
stays reachable everywhere.

## Features

- **Window config** — size, min/max, resizable, frameless, draggable, fullscreen, centered,
  always-on-top, plus runtime actions (minimize, maximize, center, raise, close).
- **Elements** — `Box`, `Row`, `Stack`, `Text`, `Divider`, `Raw`, `Responsive`.
- **Components** — Button, Progress, Spinner, Tabs, Drawer, Dialog, Scroll, Title, Image,
  a right-click context Menu with window-bounds clamping, and the form controls grouped
  under `Input`: Text, Checkbox, Slider, Select.
- **CSS-like styling** — 40+ properties, each defined in its own file, extensible by design.
- **Routing** — `NewRouter`, `Push`/`Navigate`/`Back`, `Link`.
- **Breakpoints** — `Responsive` switches views when the available width crosses a threshold.
- **Icons** — Material set, SVG path data, IconVG, plain images.
- **Component kit** — create your own components with the same helpers the built-ins use.
- **Low-level escape hatches** — `Raw`, `RunRaw`, `App.Window()`.

## Documentation

- [Getting started](getting-started.md)
- [Window configuration](window.md)
- [Elements](elements.md)
- [Style properties](styles.md)
- [Components](components/)
- [Routing](routing.md)
- [Extending the library](extending.md)
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
| `components/kit` | Helpers for custom components |
| `router` | Navigation |
| `fonts` | Font/shaper registry |

## About this documentation

The pages you are reading are Markdown files in this folder, rendered by a Vite site in
the `web/` directory. See [README.md](README.md) for the structure and how to add pages.

