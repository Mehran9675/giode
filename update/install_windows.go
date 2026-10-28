//go:build windows

package update

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func install(target, download, kind string) error {
	switch kind {
	case "exe":
		// Windows lets a running program be renamed, not overwritten.
		if err := replaceFile(target, download); err != nil {
			return err
		}
		return restart(target)
	case "installer":
		setup := download
		if !strings.HasSuffix(strings.ToLower(setup), ".exe") {
			setup = download + ".exe"
			if err := os.Rename(download, setup); err != nil {
				return fmt.Errorf("update: %w", err)
			}
		}
		// Once this process exits: install silently, then start the
		// app again.
		q := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }
		script := "Wait-Process -Id " + strconv.Itoa(os.Getpid()) + " -ErrorAction SilentlyContinue; " +
			"Start-Process -FilePath " + q(setup) + " -ArgumentList '/S' -Wait; " +
			"Remove-Item -LiteralPath " + q(setup) + " -ErrorAction SilentlyContinue; " +
			"Start-Process -FilePath " + q(target)
		cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden", "-Command", script)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000} // CREATE_NO_WINDOW
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("update: starting the installer: %w", err)
		}
		exit(0)
		return nil
	}
	return fmt.Errorf("update: can't install a %q download on Windows", kind)
}
