//go:build android || ios || js || (!linux && !darwin && !windows)

package giode

func isElevated() bool { return false }

func startElevated(a *App, exe string, args []string) error { return ErrNotSupported }
