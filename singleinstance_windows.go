//go:build windows

package giode

import "golang.org/x/sys/windows"

// stillActive is the STATUS_PENDING exit code of a running process.
const stillActive = 259

// processAlive reports whether the process with the given pid is
// running.
func processAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	defer windows.CloseHandle(h)
	var code uint32
	if err := windows.GetExitCodeProcess(h, &code); err != nil {
		return true
	}
	return code == stillActive
}
