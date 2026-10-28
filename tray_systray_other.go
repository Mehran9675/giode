//go:build windows || (linux && !android)

package giode

import (
	"runtime"

	"fyne.io/systray"
)

// startSystray runs the tray's event loop on a goroutine of its own:
// on Windows its window's messages arrive on the thread that created
// it, so the loop stays on that thread.
func startSystray(ready func()) {
	go func() {
		runtime.LockOSThread()
		systray.Run(ready, nil)
	}()
}

func trayRemove() {
	systray.Quit()
}
