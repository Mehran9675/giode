//go:build !windows

package giode

import "gioui.org/app"

// platformView handles the native window once it exists. On platforms
// other than Windows, Gio exposes no runtime window icon API (the icon
// comes from the application bundle there); file drop is set up per
// platform (fileDropView).
func (a *App) platformView(ev app.ViewEvent) {
	if a.fileDragWanted() && ev.Valid() {
		a.fileDropView(ev)
	}
}

// warnNoFileDrop logs, once, that file drops won't arrive.
func (a *App) warnNoFileDrop(why string) {
	if !a.dropWarned {
		a.dropWarned = true
		a.logf(LogLevelWarning, "giode: Config.OnFileDrop and OnFileDrag %s", why)
	}
}
