//go:build (linux && !android) || freebsd || openbsd || netbsd

package i18n

import "os"

func systemLocale() string {
	for _, v := range []string{"LC_ALL", "LC_MESSAGES", "LANG", "LANGUAGE"} {
		if l := os.Getenv(v); l != "" && l != "C" && l != "POSIX" && l != "C.UTF-8" {
			if v == "LANGUAGE" {
				// A colon-separated preference list.
				for i := 0; i < len(l); i++ {
					if l[i] == ':' {
						return l[:i]
					}
				}
			}
			return l
		}
	}
	return ""
}
