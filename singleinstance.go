package giode

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// acquireInstanceLock claims a single-instance lock for the given
// application id. It reports whether another live instance already
// holds the lock.
func acquireInstanceLock(id string) (second bool, path string, owner bool) {
	if id == "" {
		id = "giode"
	}
	path = filepath.Join(os.TempDir(), "giode-"+sanitizeID(id)+".lock")

	// The lock owner writes its pid; a stale file (dead process or
	// one that crashed before writing) is reclaimed.
	for attempt := 0; attempt < 5; attempt++ {
		data, err := os.ReadFile(path)
		if err != nil {
			break
		}
		pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err == nil && pid > 0 && processAlive(pid) {
			return true, path, false
		}
		// Empty or dead: give a racing writer a moment to fill in its
		// pid before reclaiming.
		if attempt < 4 {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		os.Remove(path)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		// Lost the race with another instance.
		return true, path, false
	}
	fmt.Fprintf(f, "%d", os.Getpid())
	f.Close()
	return false, path, true
}

// releaseInstanceLock removes the lock file of the owning instance.
func releaseInstanceLock(path string) {
	if path != "" {
		os.Remove(path)
	}
}

func sanitizeID(id string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			return r
		default:
			return '-'
		}
	}, id)
}
