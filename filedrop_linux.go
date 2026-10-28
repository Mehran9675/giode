//go:build linux && !android

package giode

import "gioui.org/app"

// File drops on Linux come through giode's Gio patches (XDND on X11,
// wl_data_device on Wayland): Gio's own event loop reads the
// protocols' messages.

var linuxFileDrag bool

func installFileDrag(a *App) {
	if a.fileDragWanted() {
		linuxFileDrag = installLinuxFileDrag(a)
	}
}

func (a *App) fileDropView(ev app.ViewEvent) {
	if !linuxFileDrag {
		a.warnNoFileDrop("need giode's Gio patches on Linux: build and run the app with an up to date giode CLI (giode upgrade)")
	}
}
