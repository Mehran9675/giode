package giode

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestParseContextAction(t *testing.T) {
	id, paths, ok := parseContextAction([]string{"-v", "--giode-action=ask", "/a b/c.txt", "/d"})
	if !ok || id != "ask" || !reflect.DeepEqual(paths, []string{"/a b/c.txt", "/d"}) {
		t.Errorf("got %q %q %v", id, paths, ok)
	}
	if _, _, ok := parseContextAction([]string{"file.txt"}); ok {
		t.Error("plain arguments parsed as an action")
	}
}

func TestContextActions_BuildAndConfig(t *testing.T) {
	b, _ := json.Marshal([]ContextAction{{ID: "ask", Label: "Ask (toml)"}, {ID: "open", Label: "Open"}})
	old := buildContextActions
	buildContextActions = base64.URLEncoding.EncodeToString(b)
	defer func() { buildContextActions = old }()
	a := &App{cfg: Config{ContextActions: []ContextAction{{ID: "ask", Label: "Ask (Go)"}, {ID: "sum", Label: "Summarize", Folders: true}}}}
	got := a.ContextActions()
	want := []ContextAction{{ID: "ask", Label: "Ask (Go)"}, {ID: "open", Label: "Open"}, {ID: "sum", Label: "Summarize", Folders: true}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v\nwant %+v", got, want)
	}
}

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
