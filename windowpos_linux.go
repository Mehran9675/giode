//go:build linux && !android

package giode

import "gioui.org/app"

// Window placement on Linux works on X11. Wayland gives apps no
// positions: the compositor places windows.

func (a *App) platformFrame() (Rect, error) {
	if _, ok := a.NativeView().(app.X11ViewEvent); ok {
		return x11Frame(a.NativeView(), a.scale())
	}
	return Rect{}, a.placementErr()
}

func (a *App) platformSetPosition(x, y float32) error {
	if _, ok := a.NativeView().(app.X11ViewEvent); ok {
		return x11SetPosition(a.NativeView(), x, y, a.scale())
	}
	return a.placementErr()
}

func (a *App) platformScreen() (Screen, error) {
	if _, ok := a.NativeView().(app.X11ViewEvent); ok {
		return x11Screen(a.NativeView(), a.scale())
	}
	return Screen{}, a.placementErr()
}

func (a *App) platformSetTopMost(on bool) bool {
	if _, ok := a.NativeView().(app.X11ViewEvent); ok {
		x11SetAbove(a.NativeView(), on)
		return true
	}
	// Gio's option, a no-op on Wayland.
	return false
}

// placementView keeps the window above others on X11 when
// Config.AlwaysOnTop asks: Gio's option doesn't reach X11 windows.
func (a *App) placementView(ev app.ViewEvent) {
	if _, ok := ev.(app.X11ViewEvent); ok && ev.Valid() && a.topMost.Load() {
		x11SetAbove(ev, true)
	}
}

func (a *App) placementErr() error {
	if a.NativeView() == nil {
		return errNoWindowYet
	}
	return ErrNotSupported
}
