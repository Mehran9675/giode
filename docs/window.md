# Window configuration

`giode.New(giode.Config{...})` creates the app. The options mirror the common window
settings of desktop app frameworks (Tauri, Wails); unsupported platform features are
ignored.

## Config fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `Title` | `string` | `""` | Window title. |
| `Size` | `image.Point` | 800x600 | Initial window size in dp. |
| `MinSize` | `image.Point` | — | Lower size bound when resizable. |
| `MaxSize` | `image.Point` | — | Upper size bound when resizable. |
| `Resizable` | `bool` | `false` | Allow user resizing. When false the size is locked to `Size`. |
| `Frameless` | `bool` | `false` | Remove the OS title bar and borders. |
| `Draggable` | `bool` | `false` | Drag the window from empty (non-widget) space. Frameless only. |
| `Background` | `color.NRGBA` | black | Window background color. |
| `Fullscreen` | `bool` | `false` | Start fullscreen. |
| `Maximized` | `bool` | `false` | Start maximized. |
| `Minimized` | `bool` | `false` | Start minimized. |
| `Center` | `bool` | `false` | Start centered on the screen. |
| `AlwaysOnTop` | `bool` | `false` | Keep above other windows. |

## Runtime window actions

```go
app.SetFullscreen(true) // or false to leave fullscreen
app.Minimize()
app.Maximize()
app.Unmaximize()
app.Center()
app.Raise()
app.Close()
```

`app.Window()` returns the underlying `*app.Window` for anything not covered here.

## Draggable windows

`Draggable: true` registers the window background as a native move area only while the
pointer hovers empty space, so widgets keep receiving clicks. The drag is OS-native: the
window moves with the mouse like a title-bar drag.
