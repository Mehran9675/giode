//go:build giode_gio2 && linux && !android

package giode

import (
	"gioui.org/app"
	"gioui.org/f32"
)

func installLinuxFileDrag(a *App) bool {
	app.GiodeSetFileDrag(func(kind int, x, y float32, uris string) {
		pos := f32.Pt(x, y)
		switch kind {
		case app.GiodeDragEnter:
			a.deliverFileDrag(FileDragEnter, pos)
		case app.GiodeDragOver:
			a.deliverFileDrag(FileDragOver, pos)
		case app.GiodeDragLeave:
			a.deliverFileDrag(FileDragLeave, pos)
		case app.GiodeDragDrop:
			a.deliverFileDrop(parseURIList(uris), pos)
		}
	})
	return true
}
