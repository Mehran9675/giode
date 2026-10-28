package giode

import (
	"image"
	"image/color"
	"net/url"
	"os"
	"sync"
	"sync/atomic"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/unit"
	"gioui.org/x/explorer"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/mehran9675/giode/components/menu"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/state"
)

// Config configures the application window. The options cover the
// common window settings of desktop applications; unsupported
// platform features are simply ignored.
type Config struct {
	// Title is the window title.
	Title string
	// Size is the initial window size. Defaults to 800x600.
	Size image.Point
	// MinSize and MaxSize bound the window size when Resizable. They
	// combine with MinWidth/MinHeight/MaxWidth/MaxHeight below; the
	// tighter bound wins.
	MinSize image.Point
	MaxSize image.Point
	// MinWidth, MinHeight, MaxWidth and MaxHeight bound the window
	// size per axis.
	MinWidth  int
	MinHeight int
	MaxWidth  int
	MaxHeight int
	// Resizable controls whether the user can resize the window. When
	// false the size is locked to Size.
	Resizable bool
	// Frameless removes the OS window frame (title bar and borders).
	Frameless bool
	// Draggable makes the window movable by dragging empty
	// (non-widget) space anywhere in it. It only applies to frameless
	// windows. To drag by a custom title bar only, leave it off and
	// use TitleBar, which marks its own area as the drag handle.
	Draggable bool
	// Background is the window background color. Theme tokens
	// ("theme.background") follow SetMode and SetPalette.
	Background string
	// StatusBar and NavigationBar color Android's status bar and
	// navigation bar. Both default to Background and follow theme
	// changes, so the bars blend with a light or dark app. Other
	// platforms ignore them.
	StatusBar     string
	NavigationBar string
	// Orientation restricts rotation on mobile platforms. Defaults to
	// OrientationAny. Desktop platforms ignore it.
	Orientation Orientation
	// FollowSystemTheme sets the theme mode (SetMode) from the
	// platform's light or dark appearance at startup and whenever it
	// changes. See SystemDarkMode for the platforms that report it.
	FollowSystemTheme bool
	// OnSystemThemeChange runs on the UI goroutine when the platform's
	// appearance changes between light and dark, with or without
	// FollowSystemTheme.
	OnSystemThemeChange func(dark bool)
	// EdgeToEdge lays the view out over the whole window, under the
	// status bar, navigation bar and cutouts, instead of inside the
	// safe area. Pad content with SafeArea where it must stay clear of
	// them. Desktop windows have no insets, so it changes nothing there.
	EdgeToEdge bool
	// Icon is the window icon, PNG or JPEG bytes. It is applied at
	// runtime where the platform supports it (Windows). The icon
	// shown by file managers comes from the executable and is set at
	// build time.
	Icon []byte
	// WindowStartState selects the initial window state. It takes
	// precedence over the Fullscreen, Maximized and Minimized flags.
	WindowStartState WindowStartState
	// Fullscreen starts the window in fullscreen mode.
	Fullscreen bool
	// Maximized starts the window maximized.
	Maximized bool
	// Minimized starts the window minimized.
	Minimized bool
	// Center starts the window centered on the screen.
	Center bool
	// AlwaysOnTop keeps the window above other windows.
	AlwaysOnTop bool
	// Logger receives Giode diagnostics at or above LogLevel.
	Logger Logger
	// LogLevel filters Logger output. Defaults to LogLevelDebug.
	LogLevel LogLevel
	// OnStartup runs before the event loop starts.
	OnStartup func()
	// OnReady runs after the first frame is presented to the screen.
	OnReady func()
	// OnShutdown runs after the window closes. On Android and iOS the
	// system may end a backgrounded app without warning, so save state
	// in OnBackground instead.
	OnShutdown func()
	// OnBackground runs when the app leaves the screen: on Android when
	// its activity stops (the user switched apps, went home, or opened
	// a system picker), on desktops when the window is minimized. On
	// Android it's the last reliable moment to save: the system may
	// freeze or end the app afterwards without warning. It runs
	// synchronously on the platform's main thread (Android) or the UI
	// goroutine (desktop): keep it short, and don't wait on the UI
	// (RunOnUI) from it. On iOS it runs on UIApplication's
	// background notification.
	OnBackground func()
	// OnForeground runs when the app comes back after OnBackground,
	// on the same goroutine.
	OnForeground func()
	// OnEvent observes every window event before giode handles it, on
	// the UI goroutine - for Gio extensions that need to see them. It
	// must not call FrameEvent.Frame. OpenFile and SaveFile already
	// get the events they need.
	OnEvent func(event.Event)
	// OnBeforeClose runs when a window close is requested. It cannot
	// veto the close.
	OnBeforeClose func()
	// OnWindowStateChange runs when the window is minimized,
	// maximized, made fullscreen or restored - by the app or by the
	// user.
	OnWindowStateChange func(WindowState)
	// OnOpenURL runs on the UI goroutine with a link that opened the
	// app or reached it while running: a URL scheme it registered
	// (giode.toml's app.schemes) or one of its universal links / Android
	// App Links (app.links). With a SingleInstanceLock, a link that
	// starts a second desktop instance reaches the running one here.
	OnOpenURL func(u *url.URL)
	// ContextActions add entries to the file manager's right-click menu
	// that open the app with the chosen paths; see ContextAction. The
	// giode CLI adds the ones in giode.toml too.
	ContextActions []ContextAction
	// OnContextAction runs on the UI goroutine when the app was started
	// from one of its context actions, with the action's ID and the
	// paths it was chosen on. With a SingleInstanceLock, a context
	// action chosen while the app runs reaches the running instance
	// here.
	OnContextAction func(id string, paths []string)
	// MinimizeToTray hides the window when it's minimized, so it leaves
	// the taskbar or Dock and lives in the tray until ShowWindow (by
	// default a click on the tray icon) brings it back. It needs a tray
	// icon: set one up with App.Tray. Desktop only.
	MinimizeToTray bool
	// CloseToTray hides the window when the user closes it, instead of
	// ending the app; Quit (a tray menu item, say) ends it. Windows and
	// macOS; on Linux, where Gio can't keep a closed window, closing
	// still quits.
	CloseToTray bool
	// MenuBarApp makes a macOS app live in the menu bar: no Dock icon
	// and no app menu, as menu bar utilities do. Pair it with App.Tray.
	// Other platforms ignore it.
	MenuBarApp bool
	// SingleInstanceLock, when set, ensures only one instance of the
	// application runs.
	SingleInstanceLock *SingleInstanceLock
	// OnFileDrop runs with the paths of files dropped onto the window
	// from the file manager, on the UI goroutine. It works when the app
	// runs as administrator too. Windows only for now; elsewhere a
	// warning is logged once.
	OnFileDrop func(paths []string)
	// MCP starts a Model Context Protocol server alongside the window,
	// so AI agents can see and drive the running app: inspect the UI,
	// click, type, take screenshots, read and set named state - plus
	// any tools added with App.MCPTool. It listens on
	// http://127.0.0.1:7600/mcp unless MCPOptions say otherwise.
	// Setting the GIODE_MCP=1 environment variable enables it too,
	// without rebuilding.
	MCP bool
	// MCPOptions configures the MCP server.
	MCPOptions MCPOptions
	// ShowFrameTimes draws a small counter in the top right corner with
	// how long frames take to build and to render, and the frame rate:
	// a development aid for finding what makes scrolling or a stream
	// stutter on a device. App.FrameStats returns the same numbers, and
	// the MCP server reports them as frame_stats.
	ShowFrameTimes bool
}

// WindowStartState selects the initial state of the window.
type WindowStartState int

const (
	// WindowStartNormal starts the window in its normal state.
	WindowStartNormal WindowStartState = iota
	// WindowStartMinimized starts the window minimized.
	WindowStartMinimized
	// WindowStartMaximized starts the window maximized.
	WindowStartMaximized
	// WindowStartFullscreen starts the window fullscreen.
	WindowStartFullscreen
)

// LogLevel filters the messages passed to Logger.
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

// Logger receives Giode diagnostics. Any *log.Logger satisfies it.
type Logger interface {
	Printf(format string, v ...any)
}

// SingleInstanceLock prevents a second instance of the application
// from running. When a second instance launches, its
// OnSecondInstanceLaunch callback runs and the instance exits without
// showing a window.
type SingleInstanceLock struct {
	// UniqueID identifies the application; it must be unique per app.
	UniqueID string
	// OnSecondInstanceLaunch runs in the second instance in place of
	// the UI.
	OnSecondInstanceLaunch func()
	// OnLaunchArgs runs in the running instance, on the UI goroutine,
	// with the command-line arguments (without the program name) of a
	// second instance that started and exited - a file opened with the
	// app, a URL. The window comes to the front either way. Context
	// actions go to Config.OnContextAction instead.
	OnLaunchArgs func(args []string)
}

// App owns a window and its render loop.
type App struct {
	cfg        Config
	window     *app.Window
	winPointer windowPointer
	ctxMenu    *menu.Menu

	icon        image.Image
	iconApplied bool

	instanceLockPath string
	instanceLock     *instanceLock
	secondInstance   bool

	winState   atomic.Int32
	winFocused atomic.Bool

	clip clipboardState

	viewMu     sync.Mutex
	view       app.ViewEvent
	dropWarned bool

	// The system bar colors last applied; see updateSystemBars.
	barsApplied bool
	statusColor color.NRGBA
	navColor    color.NRGBA

	hideKeyboard atomic.Bool
	frames       frameRecorder
	themeCheck   chan struct{}

	// explorer shows the file dialogs; it sees every window event.
	explorer *explorer.Explorer

	mcp     *mcpBridge
	uiMu    sync.Mutex
	uiCalls []func()

	// The tray icon (see Tray), and whether the window is hidden to it.
	tray         *Tray
	trayOnce     sync.Once
	windowHidden atomic.Bool
	quitting     atomic.Bool
}

// New creates an App from cfg. The window is created and shown the
// first time the app runs.
func New(cfg Config) *App {
	if cfg.Size == (image.Point{}) {
		cfg.Size = image.Pt(800, 600)
	}
	if properties.CalcColor(cfg.Background).A == 0 {
		cfg.Background = "#000000"
	}
	// Per-axis size bounds combine with the point bounds; the tighter
	// bound wins.
	if cfg.MinWidth > 0 && (cfg.MinSize.X == 0 || cfg.MinSize.X < cfg.MinWidth) {
		cfg.MinSize.X = cfg.MinWidth
	}
	if cfg.MinHeight > 0 && (cfg.MinSize.Y == 0 || cfg.MinSize.Y < cfg.MinHeight) {
		cfg.MinSize.Y = cfg.MinHeight
	}
	if cfg.MaxWidth > 0 && (cfg.MaxSize.X == 0 || cfg.MaxSize.X > cfg.MaxWidth) {
		cfg.MaxSize.X = cfg.MaxWidth
	}
	if cfg.MaxHeight > 0 && (cfg.MaxSize.Y == 0 || cfg.MaxSize.Y > cfg.MaxHeight) {
		cfg.MaxSize.Y = cfg.MaxHeight
	}
	opts := []app.Option{
		app.Title(cfg.Title),
		app.Size(unit.Dp(cfg.Size.X), unit.Dp(cfg.Size.Y)),
	}
	if cfg.Resizable {
		if cfg.MinSize != (image.Point{}) {
			opts = append(opts, app.MinSize(unit.Dp(cfg.MinSize.X), unit.Dp(cfg.MinSize.Y)))
		}
		if cfg.MaxSize != (image.Point{}) {
			opts = append(opts, app.MaxSize(unit.Dp(cfg.MaxSize.X), unit.Dp(cfg.MaxSize.Y)))
		}
	} else {
		opts = append(opts,
			app.MinSize(unit.Dp(cfg.Size.X), unit.Dp(cfg.Size.Y)),
			app.MaxSize(unit.Dp(cfg.Size.X), unit.Dp(cfg.Size.Y)),
		)
	}
	opts = append(opts, app.Decorated(!cfg.Frameless))
	switch cfg.WindowStartState {
	case WindowStartMinimized:
		opts = append(opts, app.Minimized.Option())
	case WindowStartMaximized:
		opts = append(opts, app.Maximized.Option())
	case WindowStartFullscreen:
		opts = append(opts, app.Fullscreen.Option())
	default:
		// Fall back to the individual flags.
		if cfg.Fullscreen {
			opts = append(opts, app.Fullscreen.Option())
		} else if cfg.Maximized {
			opts = append(opts, app.Maximized.Option())
		} else if cfg.Minimized {
			opts = append(opts, app.Minimized.Option())
		}
	}
	if cfg.AlwaysOnTop {
		opts = append(opts, app.TopMost(true))
	}
	if cfg.Orientation != OrientationAny {
		opts = append(opts, cfg.Orientation.option())
	}
	a := &App{cfg: cfg, themeCheck: make(chan struct{}, 1)}
	currentApp.Store(a)
	a.platformInit()
	if cfg.SingleInstanceLock == nil {
		// Restarted by an update: let the old process finish first.
		waitForUpdatedInstance()
	}
	a.setIconBytes(cfg.Icon)
	if cfg.SingleInstanceLock != nil {
		a.secondInstance, a.instanceLockPath, a.instanceLock = acquireInstanceLock(cfg.SingleInstanceLock.UniqueID)
	}
	a.window = new(app.Window)
	a.window.Option(opts...)
	a.explorer = explorer.NewExplorer(a.window)
	elements.SetLinkOpener(a.openLink)
	elements.SetLongPressFeedback(func() { a.Haptic(HapticHeavy) })
	if cfg.Center {
		a.window.Perform(system.ActionCenter)
	}
	state.SetInvalidator(a.window.Invalidate)
	if cfg.Logger != nil && cfg.LogLevel <= LogLevelDebug {
		elements.SetDebugLogger(cfg.Logger.Printf)
	}
	properties.SetThemeChangeHook(a.window.Invalidate)
	if cfg.MCP || os.Getenv("GIODE_MCP") == "1" {
		a.mcp = newMCPBridge(a, cfg.MCPOptions)
	}
	return a
}

// --- Window actions --------------------------------------------------------

// SetFullscreen switches the window in or out of fullscreen mode.
func (a *App) SetFullscreen(v bool) {
	if v {
		a.window.Option(app.Fullscreen.Option())
	} else {
		a.window.Option(app.Windowed.Option())
	}
}

// Minimize minimizes the window.
func (a *App) Minimize() {
	a.window.Perform(system.ActionMinimize)
}

// Maximize maximizes the window.
func (a *App) Maximize() {
	a.window.Perform(system.ActionMaximize)
}

// Unmaximize restores the window from its maximized state.
func (a *App) Unmaximize() {
	a.window.Perform(system.ActionUnmaximize)
}

// SetSize resizes the window to w x h dp, leaving fullscreen or
// maximized state first. The size bounds from Config (MinSize,
// MaxSize, MinWidth, ...) are widened when they'd prevent it, so
// automation can shrink a desktop app to a phone's size; a window
// that isn't Resizable stays locked, at the new size.
func (a *App) SetSize(w, h int) {
	opts := []app.Option{app.Windowed.Option(), app.Size(unit.Dp(w), unit.Dp(h))}
	if !a.cfg.Resizable {
		a.cfg.MinSize, a.cfg.MaxSize = image.Pt(w, h), image.Pt(w, h)
	} else {
		if a.cfg.MinSize.X > w || a.cfg.MinSize.Y > h {
			a.cfg.MinSize = image.Pt(min(a.cfg.MinSize.X, w), min(a.cfg.MinSize.Y, h))
		}
		if a.cfg.MaxSize != (image.Point{}) && (a.cfg.MaxSize.X < w || a.cfg.MaxSize.Y < h) {
			a.cfg.MaxSize = image.Pt(max(a.cfg.MaxSize.X, w), max(a.cfg.MaxSize.Y, h))
		}
	}
	if a.cfg.MinSize != (image.Point{}) {
		opts = append(opts, app.MinSize(unit.Dp(a.cfg.MinSize.X), unit.Dp(a.cfg.MinSize.Y)))
	}
	if a.cfg.MaxSize != (image.Point{}) {
		opts = append(opts, app.MaxSize(unit.Dp(a.cfg.MaxSize.X), unit.Dp(a.cfg.MaxSize.Y)))
	}
	a.window.Option(opts...)
}

// Center centers the window on the screen.
func (a *App) Center() {
	a.window.Perform(system.ActionCenter)
}

// Raise brings the window to the front.
func (a *App) Raise() {
	a.window.Perform(system.ActionRaise)
}

// Close requests the window to close, ending the app - also with
// CloseToTray, which only applies to the user closing it.
func (a *App) Close() {
	a.quitting.Store(true)
	a.window.Perform(system.ActionClose)
}

// Window exposes the underlying *app.Window for low-level
// customization.
func (a *App) Window() *app.Window {
	return a.window
}
