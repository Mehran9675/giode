# Window configuration

`giode.New(giode.Config{...})` creates the app and its window. Options a platform can't honor
are ignored.

## Config

| Field | Type | Default | Description |
| --- | --- | --- | --- |
| `Title` | `string` | `""` | Window title. |
| `Size` | `image.Point` | 800×600 | Initial size, in dp. |
| `Resizable` | `bool` | `false` | Allow resizing. When false, the size is fixed at `Size`. |
| `MinSize`, `MaxSize` | `image.Point` | none | Size limits for resizable windows. |
| `MinWidth`, `MinHeight`, `MaxWidth`, `MaxHeight` | `int` | none | Per-axis limits. Combined with `MinSize`/`MaxSize`, the tighter one wins. |
| `Background` | `string` | `"#000000"` | Window background, any [color](styles.md#colors). |
| `Frameless` | `bool` | `false` | Remove the OS title bar and borders. Pair with [TitleBar](components/titlebar.md). |
| `Draggable` | `bool` | `false` | Frameless only: dragging empty space anywhere moves the window. |
| `WindowStartState` | `WindowStartState` | `WindowStartNormal` | `WindowStartMinimized`, `WindowStartMaximized` or `WindowStartFullscreen`. |
| `Fullscreen`, `Maximized`, `Minimized` | `bool` | `false` | Older way to pick the start state; `WindowStartState` wins when set. |
| `Center` | `bool` | `false` | Center on screen at startup. |
| `AlwaysOnTop` | `bool` | `false` | Keep above other windows. |
| `Icon` | `[]byte` | `nil` | Window icon as PNG or JPEG bytes. Applied at runtime on Windows only. |
| `OnStartup` | `func()` | `nil` | Before the event loop starts. |
| `OnReady` | `func()` | `nil` | After the first frame is on screen. |
| `OnBeforeClose` | `func()` | `nil` | When the window is closing. It can't cancel the close. |
| `OnShutdown` | `func()` | `nil` | After the window has closed. |
| `OnWindowStateChange` | `func(WindowState)` | `nil` | After the window is minimized, maximized, made fullscreen or restored, by the app or the user. |
| `OnFileDrop` | `func([]string)` | `nil` | Paths of files dropped onto the window. Windows only; see [File drop](#file-drop). |
| `SingleInstanceLock` | `*SingleInstanceLock` | `nil` | Allow only one running instance; see below. |
| `Logger`, `LogLevel` | `Logger`, `LogLevel` | `nil`, debug | Where Giode's diagnostics go. Any `*log.Logger` works. |
| `MCP`, `MCPOptions` | `bool`, `MCPOptions` | off | Run the built-in [MCP server](mcp.md). |

`Run(view)` blocks until the window closes, then runs `OnShutdown` and exits the process.

## Window actions and state

```go
app.Minimize()
app.Maximize()
app.Unmaximize()
app.ToggleMaximize()     // maximize, or restore if maximized
app.SetFullscreen(true)  // false to leave fullscreen
app.Center()
app.Raise()
app.Close()

app.WindowState() // WindowNormal, WindowMinimized, WindowMaximized or WindowFullscreen
app.Maximized()
app.Fullscreen()
app.Focused()
```

The state methods report what the platform says, so they also reflect changes the user
makes (double-clicking the title bar, snapping the window). They're safe to call from any
goroutine.

## Frameless windows

`Frameless: true` removes the OS title bar. Add [TitleBar](components/titlebar.md) for a
custom one with drag area and window buttons. For a window with no title bar at all, use
`Draggable: true` instead, so dragging any empty space moves it; buttons and other
interactive components inside stay clickable.

## Single instance

```go
app := giode.New(giode.Config{
	SingleInstanceLock: &giode.SingleInstanceLock{
		UniqueID: "com.example.myapp",
		OnSecondInstanceLaunch: func() {
			fmt.Println("already running")
		},
	},
})
```

A second launch runs `OnSecondInstanceLaunch` and exits without opening a window. If the
first instance crashed, its stale lock is reclaimed.

## Clipboard

```go
app.SetClipboard("copied text")             // any goroutine; applied on the next frame
app.ReadClipboard(func(text string) { ... }) // the callback runs on the UI goroutine
```

Inside a `Layout` method, `kit.WriteClipboard(gtx, text)` copies right away.
[CopyButton](components/copybutton.md) is a ready-made copy control.

## File drop

```go
giode.Config{
	OnFileDrop: func(paths []string) {
		model.files = append(model.files, paths...) // UI goroutine: no locks needed
	},
}
```

This works on Windows, including when the app runs as administrator. It doesn't work on
other platforms yet; setting `OnFileDrop` there logs a warning. On Linux, Gio's event loop
consumes the drag-and-drop messages, so support needs changes in Gio itself. On macOS it
needs AppKit code that hasn't been written. There are no drag-enter or drag-leave events for
highlighting a drop target yet.

## Window icon

```go
//go:embed app.png
var icon []byte

app := giode.New(giode.Config{Icon: icon})
```

This sets the title-bar and taskbar icon at runtime on Windows. The icon file managers show
comes from the executable, so embed it at build time, for example with `rsrc` or
`goversioninfo` on Windows or the app bundle on macOS.

## Native access

- `app.Window()` returns Gio's `*app.Window`.
- `app.NativeHandle()` returns the native handle: the `HWND` on Windows, the `NSView` on macOS,
  the X11 window or Wayland surface on Linux. It's 0 until the window exists.
- `app.NativeView()` returns Gio's platform view event with all details.

## Not available

Gio doesn't expose these, so Giode can't either: hidden windows, per-button
close/minimize/maximize flags, transparent windows, positioning the window at given
coordinates, native application menus, and power events such as suspend and resume.
