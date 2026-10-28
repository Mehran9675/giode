//go:build darwin && !ios

package giode

/*
#cgo CFLAGS: -fmodules -fobjc-arc -x objective-c
#cgo LDFLAGS: -framework AppKit

#include <stdint.h>

void giodeEnableFileDrop(uintptr_t view);
*/
import "C"

import (
	"strings"
	"sync"

	"gioui.org/app"
	"gioui.org/f32"
)

var fileDropOnce sync.Once

func installFileDrag(a *App) {}

func (a *App) fileDropView(ev app.ViewEvent) {
	if v, ok := ev.(app.AppKitViewEvent); ok && v.View != 0 {
		fileDropOnce.Do(func() { C.giodeEnableFileDrop(C.uintptr_t(v.View)) })
	}
}

//export giodeFileDragEvent
func giodeFileDragEvent(kind C.int, x, y C.double) {
	if a := currentApp.Load(); a != nil {
		a.deliverFileDrag(FileDragKind(kind), f32.Pt(float32(x), float32(y)))
	}
}

//export giodeFileDrop
func giodeFileDrop(paths *C.char, x, y C.double) {
	if a := currentApp.Load(); a != nil {
		list := strings.Split(strings.TrimSuffix(C.GoString(paths), "\n"), "\n")
		if len(list) == 1 && list[0] == "" {
			list = nil
		}
		a.deliverFileDrop(list, f32.Pt(float32(x), float32(y)))
	}
}
