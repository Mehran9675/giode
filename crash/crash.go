// Package crash records the app's crashes and hands them back on the
// next start, to show the user or send to you.
//
//	func main() {
//		crash.Init(crash.Options{Version: "1.4.0"})
//		for _, r := range crash.Reports() {
//			// Ask the user, then:
//			if r.Send(ctx, "https://example.com/crashes") == nil {
//				r.Delete()
//			}
//		}
//		...
//	}
//
// A crash is an unrecovered panic or a fatal error (out of memory, a
// concurrent map write, a segmentation fault in cgo code) on any
// goroutine: the Go runtime writes its message and every goroutine's
// stack to the crash file as the app dies, so even a crash on a
// background goroutine is caught. Crashes in other languages' threads
// (Java on Android, Objective-C on Apple platforms) aren't Go's and
// aren't recorded.
package crash

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"

	"gioui.org/app"
)

// Options configure crash recording.
type Options struct {
	// Dir holds the crash files; a crashes directory in the app's data
	// directory by default.
	Dir string
	// Version is the app's version, recorded with each crash.
	Version string
	// Keep is how many reports to keep, the newest; 10 by default.
	Keep int
}

// Report is a crash from an earlier run.
type Report struct {
	// ID identifies the report; Time is when the crashed run started.
	ID   string    `json:"id"`
	Time time.Time `json:"time"`
	// Version is the app's version, GoVersion, OS and Arch the build's.
	Version   string `json:"version,omitempty"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	// Message is the panic or fatal error message; Trace the full
	// output, with every goroutine's stack.
	Message string `json:"message"`
	Trace   string `json:"trace"`

	path string
}

var (
	mu   sync.Mutex
	dir  string
	cur  *os.File
	opts Options
)

// header starts the crash file of a run; the runtime appends the
// crash after it.
type header struct {
	Time    time.Time `json:"time"`
	Version string    `json:"version,omitempty"`
}

const headerPrefix = "giode-crash "

// Init starts recording crashes for this run, and turns the previous
// run's crash, if it had one, into a report. Call it early in main.
func Init(o Options) error {
	if err := start(o); err != nil {
		return err
	}
	prune()
	return nil
}

func start(o Options) error {
	mu.Lock()
	defer mu.Unlock()
	if o.Keep <= 0 {
		o.Keep = 10
	}
	opts = o
	d := o.Dir
	if d == "" {
		data, err := app.DataDir()
		if err != nil {
			return err
		}
		d = filepath.Join(data, appName(), "crashes")
	}
	if err := os.MkdirAll(d, 0o700); err != nil {
		return err
	}
	dir = d
	current := filepath.Join(d, "current.log")
	if b, err := os.ReadFile(current); err == nil {
		if r, ok := parse(b); ok {
			id := r.Time.UTC().Format("20060102T150405.000000000")
			os.WriteFile(filepath.Join(d, id+".log"), b, 0o600)
		}
	}
	f, err := os.OpenFile(current, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	h, _ := json.Marshal(header{Time: time.Now(), Version: o.Version})
	fmt.Fprintf(f, "%s%s\n", headerPrefix, h)
	if err := debug.SetCrashOutput(f, debug.CrashOptions{}); err != nil {
		f.Close()
		return err
	}
	if cur != nil {
		cur.Close()
	}
	cur = f
	return nil
}

// Stop stops recording crashes and closes the crash file.
func Stop() {
	mu.Lock()
	defer mu.Unlock()
	debug.SetCrashOutput(nil, debug.CrashOptions{})
	if cur != nil {
		cur.Close()
		cur = nil
	}
}

// parse reads a crash file: the header line, then what the runtime
// wrote. A run that didn't crash has only the header.
func parse(b []byte) (Report, bool) {
	line, rest, _ := bytes.Cut(b, []byte("\n"))
	if !bytes.HasPrefix(line, []byte(headerPrefix)) {
		return Report{}, false
	}
	var h header
	json.Unmarshal(line[len(headerPrefix):], &h)
	trace := strings.TrimSpace(string(rest))
	if trace == "" {
		return Report{}, false
	}
	msg := trace
	if i := strings.Index(msg, "\n"); i >= 0 {
		msg = msg[:i]
	}
	return Report{
		ID: h.Time.UTC().Format("20060102T150405.000000000"), Time: h.Time, Version: h.Version,
		GoVersion: runtime.Version(), OS: runtime.GOOS, Arch: runtime.GOARCH,
		Message: msg, Trace: trace,
	}, true
}

// Reports returns the crash reports from earlier runs, newest first.
func Reports() []Report {
	mu.Lock()
	d := dir
	mu.Unlock()
	if d == "" {
		return nil
	}
	entries, _ := os.ReadDir(d)
	var out []Report
	for _, e := range entries {
		if e.Name() == "current.log" || !strings.HasSuffix(e.Name(), ".log") {
			continue
		}
		p := filepath.Join(d, e.Name())
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		if r, ok := parse(b); ok {
			r.path = p
			out = append(out, r)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Time.After(out[j].Time) })
	return out
}

// Delete removes the report, once it's been sent or dismissed.
func (r Report) Delete() error {
	if r.path == "" {
		return errors.New("crash: the report has no file")
	}
	return os.Remove(r.path)
}

// Send posts the report as JSON to url.
func (r Report) Send(ctx context.Context, url string) error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("crash: sending the report: %s", resp.Status)
	}
	return nil
}

// prune keeps the newest reports.
func prune() {
	mu.Lock()
	keep := opts.Keep
	mu.Unlock()
	rs := Reports()
	for _, r := range rs[min(len(rs), keep):] {
		r.Delete()
	}
}

// appName names the app's directory: its app ID, or the executable.
func appName() string {
	if app.ID != "" {
		return app.ID
	}
	exe, _ := os.Executable()
	return strings.TrimSuffix(filepath.Base(exe), filepath.Ext(exe))
}
