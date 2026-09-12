# Window configuration

`giode.New(giode.Config{...})` creates the app. The options mirror the common window
settings of desktop app frameworks (Tauri, Wails); unsupported platform features are
ignored.

## Config fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `Title` | `string` | `""` | Window title. |
| `Size` | `image.Point` | 800x600 | Initial window size in dp. |
| `MinSize` / `MaxSize` | `image.Point` | — | Size bounds when resizable. |
| `MinWidth` / `MinHeight` / `MaxWidth` / `MaxHeight` | `int` | — | Per-axis size bounds; combine with `MinSize`/`MaxSize` (tighter bound wins). |
| `Resizable` | `bool` | `false` | Allow user resizing. When false the size is locked to `Size`. |
| `Frameless` | `bool` | `false` | Remove the OS title bar and borders. |
| `Draggable` | `bool` | `false` | Drag the window from empty (non-widget) space. Frameless only. |
| `Background` | `string` | `"#000000"` | Window background color (any CSS hex). |
| `WindowStartState` | `WindowStartState` | `WindowStartNormal` | Initial state: `WindowStartMinimized`, `WindowStartMaximized`, `WindowStartFullscreen`. Takes precedence over the flags below. |
| `Fullscreen` / `Maximized` / `Minimized` | `bool` | `false` | Start in the given state (used when `WindowStartState` is `WindowStartNormal`). |
| `Center` | `bool` | `false` | Start centered on the screen. |
| `AlwaysOnTop` | `bool` | `false` | Keep above other windows. |
| `Logger` | `Logger` | `nil` | Receives diagnostics; any `*log.Logger` satisfies it. |
| `LogLevel` | `LogLevel` | `LogLevelDebug` | `LogLevelDebug`, `LogLevelInfo`, `LogLevelWarning`, `LogLevelError`. |
| `OnStartup` | `func()` | `nil` | Runs before the event loop starts. |
| `OnReady` | `func()` | `nil` | Runs after the first frame is presented (the Wails `OnDomReady` equivalent). |
| `OnShutdown` | `func()` | `nil` | Runs after the window closes. |
| `OnBeforeClose` | `func()` | `nil` | Runs when a close is requested. Unlike Wails it cannot veto the close. |
| `SingleInstanceLock` | `*SingleInstanceLock` | `nil` | Ensures only one instance runs; see below. |

## Single instance

```go
app := giode.New(giode.Config{
	SingleInstanceLock: &giode.SingleInstanceLock{
		UniqueID: "c9c8fd93-6758-4144-87d1-34bdb0a8bd60",
		OnSecondInstanceLaunch: func() {
			fmt.Println("already running")
		},
	},
})
```

When a second instance launches, its `OnSecondInstanceLaunch` callback runs and the
instance exits without showing a window. The lock survives crashes (a dead owner's lock
is reclaimed).

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

## Wails options mapping

| Wails option | Giode |
| --- | --- |
| `Title`, `Width`/`Height`, `DisableResize`, `Frameless`, `Fullscreen`, `AlwaysOnTop`, `BackgroundColour` | `Title`, `Size`, `Resizable`, `Frameless`, `Fullscreen`, `AlwaysOnTop`, `Background` |
| `MinWidth`/`MinHeight`/`MaxWidth`/`MaxHeight` | Same names |
| `WindowStartState` | `WindowStartState` |
| `OnStartup`, `OnShutdown`, `OnBeforeClose` | Same names (`OnBeforeClose` cannot veto) |
| `OnDomReady` | `OnReady` (first frame presented) |
| `Logger`, `LogLevel` | Same names |
| `SingleInstanceLock` | Same name |
| `CSSDragProperty`/`CSSDragValue` | Replaced by `Draggable` |
| `StartHidden`, `HideWindowOnClose` | Not applicable (Gio has no hidden windows) |
| `Menu`, `AssetServer`, `Bind`, `EnumBind`, `ErrorFormatter`, `EnableDefaultContextMenu`, `EnableFraudulentWebsiteDetection`, `BindingsAllowedOrigins`, `DragAndDrop`, `Debug` | Not applicable (webview features; Gio exposes no file-drop events) |
| `Windows.*` (transparency, Mica, themes, zoom, suspend/resume, class name, …) | Not applicable (not exposed by Gio) |
| `Mac.*`, `Linux.*` | Not applicable (not exposed by Gio) |
| `LogLevelProduction` | Not applicable (no dev/release build split in plain Go) |
