//go:build !giode_gio2 && linux && !android

package giode

func installLinuxFileDrag(a *App) bool { return false }
