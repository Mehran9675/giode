//go:build android || ios || js || (!linux && !darwin && !windows)

package giode

func pickFolder(a *App, opts FolderOptions) (string, error) { return "", ErrNotSupported }
