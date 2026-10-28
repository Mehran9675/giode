//go:build !giode_gio && !giode_gio2

package giode

func installUnderlay(f func(GPUTarget), release func(), forceGL bool) error {
	return errUnderlayUnpatched
}
