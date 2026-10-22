package giode

import (
	"image"
	"os"
	"sync"
	"sync/atomic"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/mehran9675/giode/components/menu"
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
	// Background is the window background color.
	Background string
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
	// OnShutdown runs after the window closes.
	OnShutdown func()
	// OnBeforeClose runs when a window close is requested. It cannot
	// veto the close.
	OnBeforeClose func()
	// OnWindowStateChange runs when the window is minimized,
	// maximized, made fullscreen or restored - by the app or by the
	// user.
	OnWindowStateChange func(WindowState)
	// SingleInstanceLock, when set, ensures only one instance of the
	// application runs.
	SingleInstanceLock *SingleInstanceLock
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
}

// App owns a window and its render loop.
type App struct {
	cfg        Config
	window     *app.Window
	winPointer windowPointer
	ctxMenu    *menu.Menu

	icon        image.Image
	iconApplied bool

	instanceLockPath  string
	instanceLockOwner bool
	secondInstance    bool

	winState   atomic.Int32
	winFocused atomic.Bool

	clip clipboardState

	mcp     *mcpBridge
	uiMu    sync.Mutex
	uiCalls []func()
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
	a := &App{cfg: cfg}
	a.setIconBytes(cfg.Icon)
	if cfg.SingleInstanceLock != nil {
		a.secondInstance, a.instanceLockPath, a.instanceLockOwner = acquireInstanceLock(cfg.SingleInstanceLock.UniqueID)
	}
	a.window = new(app.Window)
	a.window.Option(opts...)
	if cfg.Center {
		a.window.Perform(system.ActionCenter)
	}
	state.SetInvalidator(a.window.Invalidate)
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

// Center centers the window on the screen.
func (a *App) Center() {
	a.window.Perform(system.ActionCenter)
}

// Raise brings the window to the front.
func (a *App) Raise() {
	a.window.Perform(system.ActionRaise)
}

// Close requests the window to close.
func (a *App) Close() {
	a.window.Perform(system.ActionClose)
}

// Window exposes the underlying *app.Window for low-level
// customization.
func (a *App) Window() *app.Window {
	return a.window
}
