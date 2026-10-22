# Architecture

## Repository layout

```
app.go, config.go, ...   Root package: App, Config, the frame loop, window state,
                         clipboard, file drop, MCP wiring, and reexport.go, which
                         re-exports the public API for single-import use
styles/                  The Styles struct, Merge, dp conversion (Px)
styles/properties/       One file per style property, plus color parsing and theme
elements/                Stateless building blocks and frame-level machinery:
                         Box, Text, Pressable, Memo, Keyed, portals, semantics
components/<name>/       One package per component
components/kit/          Helpers shared by components (and yours)
state/                   State[T], read tracking, named states
router/                  Navigation
mcp/                     Standalone MCP server (JSON-RPC, HTTP and stdio transports)
fonts/                   Font registration
examples/                hello, showcase, custom, mcp, titlebar, transitions, workbench
docs/, web/              This documentation and the site that renders it
```

Packages depend downward: `components` use `elements`, which use `styles` and `kit`. The root
package ties everything to a Gio window; nothing below it imports `gioui.org/app`.

## A frame

Gio asks for a frame on input, resize, or when something requested one. The app then:

1. runs functions queued with `Post` and `RunOnUI`;
2. calls the view function and lays out the returned tree;
3. lays out the context menu and draws queued portals (popups) on top;
4. runs `AfterLayout` functions;
5. hands the frame to Gio, records the semantic tree if the MCP server is on, and drops
   memos and keyed state that weren't used.

## State model

The view is rebuilt every frame, so anything that must last longer lives in component
structs created once, in `State` values, or in `Keyed` entries. Elements keep no state of
their own, except an element with hover styles, which remembers whether it's hovered.

## Adding a style property

1. Add `styles/properties/<name>.go` with the type, constants and resolution code.
2. Add the field to `Styles` (`styles/styles.go`), to `Merge` (`styles/merge.go`), and to
   `Px` (`styles/px.go`) if it's a length.
3. Use it in `elements` or the components it applies to.
4. Re-export its constants in `reexport.go`.

## Adding a component

1. Create `components/<name>/`.
2. Hold state in a struct created by `New`; build visuals from `elements`.
3. Accept `Styles` and apply defaults with `styles.Merge`.
4. Use `kit` for shading, frame requests and drag blocking, and record semantics so
   automation can see it.
5. Re-export it from `reexport.go` and add a page under `docs/components/`.

See [Writing your own components](extending.md) for details.
