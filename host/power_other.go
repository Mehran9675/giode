//go:build android || ios || js || (!windows && !darwin && !linux && !freebsd && !openbsd && !netbsd && !dragonfly)

package host

// KeepScreenOn is supported on desktops; giode.App has the mobile
// versions, which need the app's view.
func KeepScreenOn(on bool) error { return ErrNotSupported }

// KeepAwake is supported on desktops.
func KeepAwake(reason string) (release func(), err error) { return func() {}, ErrNotSupported }
