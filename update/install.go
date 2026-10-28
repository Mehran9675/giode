package update

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// WaitEnv names the environment variable an updated app is restarted
// with: the old process's pid, which giode waits for to exit before
// taking the single-instance lock.
const WaitEnv = "GIODE_UPDATE_WAIT_PID"

// restart starts exe with this process's arguments, after this process
// exits, and exits.
func restart(exe string) error {
	cmd := exec.Command(exe, os.Args[1:]...)
	cmd.Env = append(os.Environ(), WaitEnv+"="+strconv.Itoa(os.Getpid()))
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("update: restarting: %w", err)
	}
	exit(0)
	return nil
}

// exit ends the process; tests replace it.
var exit = os.Exit

// copyFile copies src to dst with mode.
func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// replaceFile puts the download in target's place: copied next to it
// (a rename can't cross file systems), the old file kept as .old.
func replaceFile(target, download string) error {
	next := target + ".new"
	if err := copyFile(download, next, 0o755); err != nil {
		return fmt.Errorf("update: %w", err)
	}
	os.RemoveAll(target + ".old")
	if err := os.Rename(target, target+".old"); err != nil {
		os.Remove(next)
		return fmt.Errorf("update: %w", err)
	}
	if err := os.Rename(next, target); err != nil {
		os.Rename(target+".old", target)
		return fmt.Errorf("update: %w", err)
	}
	os.Remove(download)
	return nil
}

// sameDirTemp makes a temporary directory beside path.
func sameDirTemp(path string) (string, error) {
	return os.MkdirTemp(filepath.Dir(path), ".giode-update-")
}
