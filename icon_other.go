//go:build !windows

package giode

import "gioui.org/app"

// platformView handles the native window once it exists. On platforms
// other than Windows, Gio exposes no runtime window icon API (the icon
// comes from the application bundle there), and file drop isn't
// available: Gio's own event loop consumes the X11 and Wayland
// drag-and-drop protocol messages, and AppKit drag support isn't
// implemented yet.
func (a *App) platformView(ev app.ViewEvent) {
	if a.cfg.OnFileDrop != nil && ev.Valid() && !a.dropWarned {
		a.dropWarned = true
		a.logf(LogLevelWarning, "giode: Config.OnFileDrop is only supported on Windows")
	}
}
