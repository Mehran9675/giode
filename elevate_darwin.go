//go:build darwin && !ios

package giode

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func isElevated() bool { return os.Geteuid() == 0 }

// startElevated asks for an administrator password through AppleScript
// and starts the app as root in the background.
func startElevated(a *App, exe string, args []string) error {
	parts := []string{shellQuote(exe)}
	for _, s := range args {
		parts = append(parts, shellQuote(s))
	}
	script := "do shell script " + appleString(strings.Join(parts, " ")+" >/dev/null 2>&1 &") + " with administrator privileges"
	out, err := exec.Command("osascript", "-e", script).CombinedOutput()
	var ee *exec.ExitError
	if errors.As(err, &ee) && strings.Contains(string(out), "-128") {
		return ErrDeclined
	}
	if err != nil {
		return errors.New("giode: " + strings.TrimSpace(string(out)))
	}
	return nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func appleString(s string) string {
	return `"` + strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(s) + `"`
}
