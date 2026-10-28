//go:build giode_gio && !giode_gio2

package giode

import (
	"errors"

	"gioui.org/app"
)

// installUnderlay hands the underlay to giode's patched Gio, from a
// giode CLI without the giode_gio2 patches: no release or after-frame
// hooks.
func installUnderlay(f func(GPUTarget), release func(), forceGL bool) error {
	if release != nil {
		return errors.New("giode: Config.UnderlayRelease needs newer Gio patches: update the giode CLI (giode upgrade)")
	}
	var hook func(app.GiodeTarget)
	if f != nil {
		hook = func(t app.GiodeTarget) {
			f(gpuTarget(t.API, t.Framebuffer, t.Handle, t.Width, t.Height))
		}
	}
	app.GiodeSetUnderlay(hook, forceGL)
	return nil
}
