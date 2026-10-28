package host

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestSingleInstance_ForwardsLaunch(t *testing.T) {
	id := "giode-test-" + filepath.Base(t.TempDir())
	second, path, owner := acquireInstanceLock(id)
	if second || owner == nil {
		t.Fatal("first instance didn't get the lock")
	}
	defer releaseInstanceLock(path, owner)
	// Windows has no permission bits; its temporary directory is the
	// user's own.
	if info, err := os.Stat(path); runtime.GOOS != "windows" && (err != nil || info.Mode().Perm()&0o077 != 0) {
		t.Errorf("lock file mode %v (%v), want private", info.Mode(), err)
	}
	got := make(chan []string, 1)
	owner.serveLaunches(func(args []string) { got <- args })

	second, path2, owner2 := acquireInstanceLock(id)
	if !second || owner2 != nil || path2 != path {
		t.Fatal("second instance got the lock")
	}
	if err := forwardLaunch(path2, []string{"--giode-action=ask", "/tmp/x.txt"}); err != nil {
		t.Fatal(err)
	}
	select {
	case args := <-got:
		if !reflect.DeepEqual(args, []string{"--giode-action=ask", "/tmp/x.txt"}) {
			t.Errorf("forwarded %q", args)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("the launch never arrived")
	}

	// A wrong token is refused.
	data, _ := os.ReadFile(path)
	pid, port, _ := parseLock(data)
	bad := filepath.Join(t.TempDir(), "bad.lock")
	os.WriteFile(bad, []byte(fmtLock(pid, port, "nope")), 0o600)
	if err := forwardLaunch(bad, nil); err == nil {
		t.Error("a launch with the wrong token was accepted")
	}
}

func fmtLock(pid, port int, token string) string {
	return strconv.Itoa(pid) + "\n" + strconv.Itoa(port) + "\n" + token + "\n"
}

// TestLockInstance checks the exported API: the first instance holds
// the lock and serves, the second forwards.
func TestLockInstance(t *testing.T) {
	id := "giode-test-api-" + filepath.Base(t.TempDir())
	first := LockInstance(id)
	if first.Running() {
		t.Fatal("first instance sees another running")
	}
	defer first.Release()
	got := make(chan []string, 1)
	first.Serve(func(args []string) { got <- args })
	second := LockInstance(id)
	if !second.Running() {
		t.Fatal("second instance got the lock")
	}
	second.Release() // must not free the first's lock
	if err := second.Forward([]string{"movie.mkv"}); err != nil {
		t.Fatal(err)
	}
	select {
	case a := <-got:
		if len(a) != 1 || a[0] != "movie.mkv" {
			t.Errorf("%q", a)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("no launch")
	}
}
