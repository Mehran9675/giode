//go:build linux && !android

package giode

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func isElevated() bool { return os.Geteuid() == 0 }

// startElevated runs the app through pkexec, passing on what a GUI
// needs from the environment, which pkexec clears. pkexec asks for the
// password and then execs the command in its own process, so the app
// has started once that process runs exe; pkexec exiting first means
// the user declined or polkit refused.
func startElevated(a *App, exe string, args []string) error {
	if _, err := exec.LookPath("pkexec"); err != nil {
		return fmt.Errorf("giode: pkexec isn't installed: %w", ErrNotSupported)
	}
	cmdArgs := []string{"env"}
	for _, k := range []string{"DISPLAY", "XAUTHORITY", "WAYLAND_DISPLAY", "XDG_RUNTIME_DIR", "XDG_SESSION_TYPE", "DBUS_SESSION_BUS_ADDRESS", "LANG", "HOME"} {
		if v, ok := os.LookupEnv(k); ok {
			cmdArgs = append(cmdArgs, k+"="+v)
		}
	}
	cmdArgs = append(append(cmdArgs, exe), args...)
	cmd := exec.Command("pkexec", cmdArgs...)
	if err := cmd.Start(); err != nil {
		return err
	}
	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()
	cmdline := "/proc/" + strconv.Itoa(cmd.Process.Pid) + "/cmdline"
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case err := <-exited:
			var ee *exec.ExitError
			if errors.As(err, &ee) && (ee.ExitCode() == 126 || ee.ExitCode() == 127) {
				return ErrDeclined
			}
			if err == nil {
				// The elevated copy already ran and quit.
				return nil
			}
			return err
		case <-tick.C:
			if b, err := os.ReadFile(cmdline); err == nil && bytes.HasPrefix(b, append([]byte(exe), 0)) {
				return nil
			}
		}
	}
}
