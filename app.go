package giode

import (
	"fmt"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"

	"github.com/mehran9675/giode/components/menu"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles/properties"
)

// SetContextMenu registers a menu to open with a right-click anywhere
// in the window. The app registers the right-click area and renders
// the menu on top automatically.
func (a *App) SetContextMenu(m *menu.Menu) {
	a.ctxMenu = m
}

// Run drives the event loop and renders view every frame until the
// window closes. view runs on every frame: build the UI from scratch,
// but reuse stateful components created before Run.
//
// Run blocks for the lifetime of the application: when the window
// closes it releases the single-instance lock, runs OnShutdown and
// terminates the process. Any window close error is printed to
// stderr.
func (a *App) Run(view func() Element) error {
	if a.secondInstance {
		a.logf(LogLevelInfo, "giode: second instance exiting")
		if a.cfg.SingleInstanceLock.OnSecondInstanceLaunch != nil {
			a.cfg.SingleInstanceLock.OnSecondInstanceLaunch()
		}
		return nil
	}
	if a.cfg.OnStartup != nil {
		a.cfg.OnStartup()
	}
	a.logf(LogLevelInfo, "giode: app started")
	go func() {
		err := a.runLoop(view)
		if a.instanceLockOwner {
			releaseInstanceLock(a.instanceLockPath)
		}
		if a.cfg.OnShutdown != nil {
			a.cfg.OnShutdown()
		}
		a.logf(LogLevelDebug, "giode: app stopped")
		if err != nil {
			fmt.Fprintln(os.Stderr, "giode:", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	// app.Main blocks forever on most platforms; the process exits
	// from the window goroutine above once the window closes.
	app.Main()
	return nil
}

// RunRaw is Run with a plain Gio layout function.
func (a *App) RunRaw(view func(gtx layout.Context) layout.Dimensions) error {
	return a.Run(func() Element { return elements.Raw(view) })
}

func (a *App) runLoop(view func() Element) error {
	var ops op.Ops
	ready := false
	for {
		switch e := a.window.Event().(type) {
		case app.DestroyEvent:
			if a.cfg.OnBeforeClose != nil {
				a.cfg.OnBeforeClose()
			}
			a.logf(LogLevelDebug, "giode: window closed")
			return e.Err
		case app.ViewEvent:
			a.handleViewEvent(e)
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			paint.Fill(gtx.Ops, properties.CalcColor(a.cfg.Background))
			if a.cfg.Frameless && a.cfg.Draggable {
				a.drag.layout(gtx)
			}
			if a.ctxMenu != nil {
				a.ctxMenu.Context(gtx)
			}
			view().Layout(gtx)
			if a.ctxMenu != nil {
				a.ctxMenu.Layout(gtx)
			}
			e.Frame(gtx.Ops)
			if !ready {
				ready = true
				a.logf(LogLevelDebug, "giode: first frame presented")
				if a.cfg.OnReady != nil {
					a.cfg.OnReady()
				}
			}
		}
	}
}

// logf forwards a message to the configured logger when the level is
// enabled.
func (a *App) logf(level LogLevel, format string, v ...any) {
	if a.cfg.Logger != nil && level >= a.cfg.LogLevel {
		a.cfg.Logger.Printf(format, v...)
	}
}
