//go:build android || ios || (!windows && !linux && !darwin)

package giode

import "gioui.org/app"

func installFileDrag(a *App) {}

func (a *App) fileDropView(ev app.ViewEvent) {
	a.warnNoFileDrop("aren't supported on this platform")
}
