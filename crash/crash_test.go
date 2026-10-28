package crash

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestMain runs the crashing child when asked to.
func TestMain(m *testing.M) {
	if d := os.Getenv("GIODE_CRASH_CHILD"); d != "" {
		Init(Options{Dir: d, Version: "1.2.3"})
		done := make(chan bool)
		go func() {
			var m map[string]int
			m["boom"] = 1 // panics on a background goroutine
			done <- true
		}()
		<-done
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func crashChild(t *testing.T, dir string) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=NONE")
	cmd.Env = append(os.Environ(), "GIODE_CRASH_CHILD="+dir)
	if err := cmd.Run(); err == nil {
		t.Fatal("the child didn't crash")
	}
}

func TestCrashIsReportedNextRun(t *testing.T) {
	dir := t.TempDir()
	t.Cleanup(Stop)
	crashChild(t, dir)
	if err := Init(Options{Dir: dir, Version: "1.2.4"}); err != nil {
		t.Fatal(err)
	}
	rs := Reports()
	if len(rs) != 1 {
		t.Fatalf("%d reports, want 1", len(rs))
	}
	r := rs[0]
	if !strings.Contains(r.Message, "assignment to entry in nil map") || r.Version != "1.2.3" {
		t.Errorf("report %q, version %q", r.Message, r.Version)
	}
	if !strings.Contains(r.Trace, "goroutine") || !strings.Contains(r.Trace, "crash_test.go") {
		t.Errorf("trace lacks the stack:\n%s", r.Trace)
	}

	// Sending posts the JSON.
	var got Report
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		json.NewDecoder(req.Body).Decode(&got)
	}))
	defer srv.Close()
	if err := r.Send(context.Background(), srv.URL); err != nil {
		t.Fatal(err)
	}
	if got.ID != r.ID || got.Trace == "" {
		t.Errorf("server got %+v", got)
	}
	if err := r.Delete(); err != nil {
		t.Fatal(err)
	}
	if len(Reports()) != 0 {
		t.Error("report left after Delete")
	}

	// A run that doesn't crash leaves no report.
	if err := Init(Options{Dir: dir}); err != nil {
		t.Fatal(err)
	}
	if len(Reports()) != 0 {
		t.Error("a clean run made a report")
	}
}

func TestKeepsNewest(t *testing.T) {
	dir := t.TempDir()
	t.Cleanup(Stop)
	for i := 0; i < 3; i++ {
		crashChild(t, dir)
		Init(Options{Dir: dir, Keep: 2})
		time.Sleep(10 * time.Millisecond)
	}
	if n := len(Reports()); n != 2 {
		t.Errorf("%d reports kept, want 2", n)
	}
}
