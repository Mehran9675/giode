# Window configuration

`giode.New(giode.Config{...})` creates the app. The options cover the common window
settings of desktop applications; unsupported platform features are ignored.

## Config fields

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `Title` | `string` | `""` | Window title. |
| `Size` | `image.Point` | 800x600 | Initial window size in dp. |
| `MinSize` / `MaxSize` | `image.Point` | — | Size bounds when resizable. |
| `MinWidth` / `MinHeight` / `MaxWidth` / `MaxHeight` | `int` | — | Per-axis size bounds; combine with `MinSize`/`MaxSize` (tighter bound wins). |
| `Resizable` | `bool` | `false` | Allow user resizing. When false the size is locked to `Size`. |
| `Frameless` | `bool` | `false` | Remove the OS title bar and borders. |
| `Draggable` | `bool` | `false` | Drag the window from empty (non-widget) space anywhere. Frameless only. For a custom title bar, use [TitleBar](components/titlebar.md) instead. |
| `OnFileDrop` | `func(paths []string)` | `nil` | Files dropped onto the window from the file manager, on the UI goroutine. Windows only for now. |
| `OnWindowStateChange` | `func(WindowState)` | `nil` | Runs when the window is minimized, maximized, made fullscreen or restored. |
| `Background` | `string` | `"#000000"` | Window background color (any CSS hex). |
| `Icon` | `[]byte` | `nil` | Window icon, PNG or JPEG bytes. Applied at runtime on Windows. |
| `WindowStartState` | `WindowStartState` | `WindowStartNormal` | Initial state: `WindowStartMinimized`, `WindowStartMaximized`, `WindowStartFullscreen`. Takes precedence over the flags below. |
| `Fullscreen` / `Maximized` / `Minimized` | `bool` | `false` | Start in the given state (used when `WindowStartState` is `WindowStartNormal`). |
| `Center` | `bool` | `false` | Start centered on the screen. |
| `AlwaysOnTop` | `bool` | `false` | Keep above other windows. |
| `Logger` | `Logger` | `nil` | Receives diagnostics; any `*log.Logger` satisfies it. |
| `LogLevel` | `LogLevel` | `LogLevelDebug` | `LogLevelDebug`, `LogLevelInfo`, `LogLevelWarning`, `LogLevelError`. |
| `OnStartup` | `func()` | `nil` | Runs before the event loop starts. |
| `OnReady` | `func()` | `nil` | Runs after the first frame is presented. |
| `OnShutdown` | `func()` | `nil` | Runs after the window closes. |
| `OnBeforeClose` | `func()` | `nil` | Runs when a close is requested. It cannot veto the close. |
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

## Window icon

Pass PNG or JPEG bytes and the icon is applied to the window (title bar and taskbar) at
runtime on Windows:

```go
icon, _ := os.ReadFile("app.png") // or embed the bytes with go:embed
app := giode.New(giode.Config{Icon: icon})
```

The icon shown by file managers comes from the executable, not the running window, and
is embedded at build time instead:

```sh
rsrc -ico app.ico -o app.syso   # or goversioninfo / windres
go build
```

## Not supported

A few common window features are deliberately absent because Gio does not expose them:
hidden windows, closable/minimizable/maximizable flags, window transparency, positioning
by coordinates, file-drop events, native application menus and power-state callbacks
(suspend/resume).

## Clipboard

```go
app.SetClipboard("copied text")             // from any goroutine; applied on the next frame
app.ReadClipboard(func(text string) { ... }) // fn runs on the UI goroutine
```

Inside a `Layout`, `kit.WriteClipboard(gtx, text)` copies immediately. For a ready-made
control, `giode.CopyButton(func() string { return value })` copies when clicked and shows a
check mark for a moment (`Labels`, `Feedback`, `OnCopy` customize it).

## Back button

`Dialog`, `Drawer` and `Menu` close on Escape and on the system Back button (Android)
while open, and `Router` goes back while it has history. For anything else, register a
handler from `Layout`, every frame, for as long as Back has something to do:

```go
func (p *EditorPage) Layout(gtx layout.Context) layout.Dimensions {
	if p.dirty {
		giode.HandleBack(func() { p.confirmDiscard = true })
	}
	...
}
```

The handler registered last in the frame wins, so an overlay on top takes Back before the
page beneath it, and only one handler runs. Handlers run after the frame's layout, so they
may change anything the view reads. When nothing is registered, Back is left to the
platform: on Android it closes the app as usual. `HandleBack` reacts to Back only;
`elements.HandleDismiss` also reacts to Escape and is what overlays use.

Handlers are registered during layout rather than consulted when Back is pressed because
Android decides whether the app handled Back the moment it's pressed, from the handlers
registered in the last frame - a callback that returned false later could no longer hand
the press back to the system.

## File drop

```go
giode.Config{
	OnFileDrop: func(paths []string) {
		model.files = append(model.files, paths...) // UI goroutine: no locks needed
	},
}
```

On Windows the app subclasses its own window and accepts `WM_DROPFILES`, including when it
runs as administrator (it lets the drop messages through Windows' privilege filter, which
would otherwise silently block drops from a normal Explorer).

Other platforms don't support it yet, and a warning is logged once if `OnFileDrop` is set
there. Gio's own event loop consumes the X11 and Wayland drag-and-drop protocol messages, so
Linux support needs changes in Gio itself. macOS support needs AppKit drag registration,
which is not implemented yet. Drag enter/leave events (for highlighting a drop zone) aren't
available yet either; they need Windows' COM `IDropTarget` instead of `WM_DROPFILES`.

## Native window handle

`app.NativeHandle()` returns the native window handle: the `HWND` on Windows, the `NSView`
on macOS, the X11 window or `wl_surface` on Linux. `app.NativeView()` returns Gio's full
platform view event for other integrations. Both are zero until the window exists.
