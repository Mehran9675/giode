//go:build !windows

package giode

import "gioui.org/app"

// handleViewEvent applies the configured icon once the platform
// window handle is known. On platforms other than Windows the runtime
// icon API is not exposed by Gio; the window icon is set at build
// time there instead.
func (a *App) handleViewEvent(ev app.ViewEvent) {}
