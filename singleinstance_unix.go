//go:build !windows

package giode

import "syscall"

// processAlive reports whether the process with the given pid is
// running.
func processAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	err := syscall.Kill(pid, 0)
	return err == nil || err == syscall.EPERM
}
