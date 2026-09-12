package giode

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/mehran9675/giode/components/menu"
	"github.com/mehran9675/giode/elements"
)

// SetContextMenu registers a menu to open with a right-click anywhere
// in the window, including on top of other widgets. The app registers
// the right-click area (last, so it isn't blocked by the rest of the
// UI — see menu.Context) and renders the menu on top automatically.
func (a *App) SetContextMenu(m *menu.Menu) {
	a.ctxMenu = m
}

// Run drives the event loop and renders view every frame until the
// window closes. view runs on every frame: build the UI from scratch,
// but reuse stateful components created before Run.
//
// Run blocks for the lifetime of the application; on desktop it
// returns the window close error, if any.
func (a *App) Run(view func() Element) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- a.runLoop(view)
	}()
	app.Main()
	return <-errCh
}

// RunRaw is Run with a plain Gio layout function.
func (a *App) RunRaw(view func(gtx layout.Context) layout.Dimensions) error {
	return a.Run(func() Element { return elements.Raw(view) })
}

func (a *App) runLoop(view func() Element) error {
	var ops op.Ops
	color := properties.CalcColor(a.cfg.Background)
	for {
		switch e := a.window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			paint.Fill(gtx.Ops, color)
			if a.cfg.Frameless && a.cfg.Draggable {
				a.drag.layout(gtx)
			}
			view().Layout(gtx)
			if a.ctxMenu != nil {
				a.ctxMenu.Context(gtx)
				a.ctxMenu.Layout(gtx)
			}
			e.Frame(gtx.Ops)
		}
	}
}
