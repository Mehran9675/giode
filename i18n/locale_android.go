//go:build android

package i18n

import (
	"os"
	"os/exec"
	"strings"
)

// systemLocale reads the device's locale property; Android doesn't put
// it in the environment.
func systemLocale() string {
	if l := os.Getenv("LANG"); l != "" {
		return l
	}
	for _, prop := range []string{"persist.sys.locale", "ro.product.locale"} {
		if out, err := exec.Command("/system/bin/getprop", prop).Output(); err == nil {
			if l := strings.TrimSpace(string(out)); l != "" {
				return l
			}
		}
	}
	return ""
}
