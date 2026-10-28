//go:build android || ios || (!windows && !linux && !darwin)

package giode

import (
	"runtime"

	"gioui.org/app"
)

// On Android, iOS and the web the window is the screen; other systems
// aren't supported.

func mobileWindow() bool {
	return runtime.GOOS == "android" || runtime.GOOS == "ios" || runtime.GOOS == "js"
}

func (a *App) platformScreen() (Screen, error) {
	if !mobileWindow() {
		return Screen{}, ErrNotSupported
	}
	return a.windowScreen()
}

func (a *App) platformFrame() (Rect, error) {
	s, err := a.platformScreen()
	return s.Bounds, err
}

func (a *App) platformSetPosition(x, y float32) error { return ErrNotSupported }

func (a *App) platformSetTopMost(on bool) bool { return false }

func (a *App) placementView(ev app.ViewEvent) {}
