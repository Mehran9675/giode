//go:build darwin && !ios

package giode

/*
#cgo CFLAGS: -fmodules -fobjc-arc -x objective-c
#cgo LDFLAGS: -framework AppKit

#include <stdint.h>

void giodeWatchPlacement(uintptr_t view);
int giodePlacement(double *out);
void giodeSetPosition(uintptr_t view, double x, double y);
*/
import "C"

import "gioui.org/app"

// macOS measures windows in points, which are Gio's dp.

func (a *App) placement() ([13]float64, error) {
	var p [13]C.double
	if a.NativeView() == nil || C.giodePlacement(&p[0]) == 0 {
		return [13]float64{}, errNoWindowYet
	}
	var out [13]float64
	for i, v := range p {
		out[i] = float64(v)
	}
	return out, nil
}

func rectAt(p []float64) Rect {
	return Rect{X: float32(p[0]), Y: float32(p[1]), Width: float32(p[2]), Height: float32(p[3])}
}

func (a *App) platformFrame() (Rect, error) {
	p, err := a.placement()
	return rectAt(p[0:4]), err
}

func (a *App) platformScreen() (Screen, error) {
	p, err := a.placement()
	return Screen{Bounds: rectAt(p[4:8]), WorkArea: rectAt(p[8:12]), Scale: float32(p[12])}, err
}

func (a *App) platformSetPosition(x, y float32) error {
	ev, ok := a.NativeView().(app.AppKitViewEvent)
	if !ok || ev.View == 0 {
		return errNoWindowYet
	}
	C.giodeSetPosition(C.uintptr_t(ev.View), C.double(x), C.double(y))
	return nil
}

// platformSetTopMost leaves always on top to Gio.
func (a *App) platformSetTopMost(on bool) bool { return false }

func (a *App) placementView(ev app.ViewEvent) {
	if v, ok := ev.(app.AppKitViewEvent); ok && v.View != 0 {
		C.giodeWatchPlacement(C.uintptr_t(v.View))
	}
}
