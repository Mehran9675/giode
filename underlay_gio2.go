//go:build giode_gio2

package giode

import "gioui.org/app"

// installUnderlay hands the underlay and its hooks to giode's patched
// Gio.
func installUnderlay(f func(GPUTarget), release func(), forceGL bool) error {
	var hook, after func(app.GiodeTarget)
	if f != nil {
		hook = func(t app.GiodeTarget) {
			f(gpuTarget(t.API, t.Framebuffer, t.Handle, t.Width, t.Height))
		}
		after = func(t app.GiodeTarget) {
			underlayFrameDone(gpuTarget(t.API, t.Framebuffer, t.Handle, t.Width, t.Height))
		}
	}
	app.GiodeSetUnderlay(hook, forceGL)
	app.GiodeSetUnderlayHooks(release, after)
	return nil
}
