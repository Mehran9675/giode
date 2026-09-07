package giode

import (
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"

	"github.com/mehran9675/giode/components/menu"
)

// Config configures the application window. The options mirror the
// common window settings of desktop app frameworks such as Tauri and
// Wails; unsupported platform features are simply ignored.
type Config struct {
	// Title is the window title.
	Title string
	// Size is the initial window size. Defaults to 800x600.
	Size image.Point
	// MinSize and MaxSize bound the window size when Resizable.
	MinSize image.Point
	MaxSize image.Point
	// Resizable controls whether the user can resize the window. When
	// false the size is locked to Size.
	Resizable bool
	// Frameless removes the OS window frame (title bar and borders).
	Frameless bool
	// Draggable makes the window movable by dragging empty
	// (non-widget) space. It only applies to frameless windows.
	Draggable bool
	// Background is the window background color.
	Background color.NRGBA
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
}

// App owns a window and its render loop.
type App struct {
	cfg     Config
	window  *app.Window
	drag    dragState
	ctxMenu *menu.Menu
}

// New creates an App from cfg. The window is created and shown the
// first time the app runs.
func New(cfg Config) *App {
	if cfg.Size == (image.Point{}) {
		cfg.Size = image.Pt(800, 600)
	}
	if cfg.Background.A == 0 {
		cfg.Background = color.NRGBA{A: 0xff}
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
	if cfg.Fullscreen {
		opts = append(opts, app.Fullscreen.Option())
	} else if cfg.Maximized {
		opts = append(opts, app.Maximized.Option())
	} else if cfg.Minimized {
		opts = append(opts, app.Minimized.Option())
	}
	if cfg.AlwaysOnTop {
		opts = append(opts, app.TopMost(true))
	}
	a := &App{cfg: cfg}
	a.window = new(app.Window)
	a.window.Option(opts...)
	if cfg.Center {
		a.window.Perform(system.ActionCenter)
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
