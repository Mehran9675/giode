package giode

import (
	"os/exec"
	"syscall"
)

// revealFile opens Explorer with path selected.
func revealFile(path string) error {
	cmd := exec.Command("explorer.exe")
	// Explorer parses its own command line: the path must be quoted
	// after "/select," rather than as a separate, escaped argument.
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: `explorer.exe /select,"` + path + `"`}
	// Explorer exits with status 1 even when it succeeds, so it isn't
	// waited for.
	return cmd.Start()
}
