//go:build !giode_gio

package giode

func installUnderlay(f func(GPUTarget), forceGL bool) error { return errUnderlayUnpatched }
