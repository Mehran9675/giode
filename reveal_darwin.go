//go:build darwin && !ios

package giode

import "os/exec"

// revealFile shows path selected in a Finder window.
func revealFile(path string) error {
	return exec.Command("open", "-R", path).Run()
}
