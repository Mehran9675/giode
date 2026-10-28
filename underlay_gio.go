//go:build giode_gio

package giode

import (
	"image"

	"gioui.org/app"
)

// installUnderlay hands the underlay to giode's patched Gio.
func installUnderlay(f func(GPUTarget), forceGL bool) error {
	var hook func(app.GiodeTarget)
	if f != nil {
		hook = func(t app.GiodeTarget) {
			f(GPUTarget{API: GPUAPI(t.API), Framebuffer: t.Framebuffer, Texture: t.Handle, Size: image.Pt(t.Width, t.Height)})
		}
	}
	app.GiodeSetUnderlay(hook, forceGL)
	return nil
}
