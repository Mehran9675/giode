package update

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func keys(t *testing.T) (pub, priv string) {
	p, k, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(p), base64.StdEncoding.EncodeToString(k)
}

func TestNewer(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want bool
	}{
		{"1.2.10", "1.2.9", true}, {"1.3", "1.2.9", true}, {"2.0.0", "10.0.0", false},
		{"1.2.3", "1.2.3", false}, {"v1.2.4", "1.2.3", true}, {"1.2.4-beta", "1.2.3", true}, {"1.2.3", "", true},
	} {
		if got := Newer(c.a, c.b); got != c.want {
			t.Errorf("Newer(%q, %q) = %v", c.a, c.b, got)
		}
	}
}

// server serves a signed feed for platform with payload as its
// download; tamper corrupts what's served after signing.
type server struct {
	feed, sig, payload []byte
}

func newServer(t *testing.T, priv, version string, payload []byte, mutate func(*Feed)) (*server, *httptest.Server) {
	s := &server{payload: payload}
	sum := sha256.Sum256(payload)
	f := Feed{Version: version, Notes: "Faster.", Assets: map[string]Asset{
		"test-plat": {URL: "app-new", SHA256: hex.EncodeToString(sum[:]), Size: int64(len(payload)), Kind: "binary"},
	}}
	if mutate != nil {
		mutate(&f)
	}
	s.feed, _ = json.Marshal(f)
	sig, err := Sign(priv, s.feed)
	if err != nil {
		t.Fatal(err)
	}
	s.sig = []byte(sig)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/feed.json":
			w.Write(s.feed)
		case "/feed.json.sig":
			w.Write(s.sig)
		case "/app-new":
			w.Write(s.payload)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return s, srv
}

func TestCheckAndDownload(t *testing.T) {
	pub, priv := keys(t)
	exe := filepath.Join(t.TempDir(), "app")
	os.WriteFile(exe, []byte("old"), 0o755)
	s, srv := newServer(t, priv, "1.3.0", []byte("new program"), nil)
	u := &Updater{FeedURL: srv.URL + "/feed.json", PublicKey: pub, Version: "1.2.0", Platform: "test-plat", Executable: exe}
	ctx := context.Background()

	rel, err := u.Check(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if rel == nil || rel.Version != "1.3.0" || rel.Notes != "Faster." {
		t.Fatalf("release %+v", rel)
	}
	var seen int64
	dl, err := u.Download(ctx, rel, func(done, total int64) { seen = done })
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dl.Path)
	if b, _ := os.ReadFile(dl.Path); string(b) != "new program" || seen != int64(len("new program")) {
		t.Errorf("downloaded %q, progress %d", b, seen)
	}

	// Up to date.
	u.Version = "1.3.0"
	if rel, err := u.Check(ctx); err != nil || rel != nil {
		t.Errorf("up to date: %v %v", rel, err)
	}
	u.Version = "1.2.0"

	// A tampered feed fails the signature.
	good := s.feed
	s.feed = []byte(strings.Replace(string(good), "1.3.0", "9.9.9", 1))
	if _, err := u.Check(ctx); err == nil || !strings.Contains(err.Error(), "signature") {
		t.Errorf("tampered feed: %v", err)
	}
	s.feed = good

	// A tampered download fails the digest.
	s.payload = []byte("new progrAm") // same length, other bytes
	if _, err := u.Download(ctx, rel, nil); err == nil || !strings.Contains(err.Error(), "SHA-256") {
		t.Errorf("tampered download: %v", err)
	}
	s.payload = []byte("new program, much longer than the feed says")
	if _, err := u.Download(ctx, rel, nil); err == nil || !strings.Contains(err.Error(), "larger") {
		t.Errorf("oversized download: %v", err)
	}

	// Another key's signature is refused.
	other, _ := keys(t)
	u.PublicKey = other
	if _, err := u.Check(ctx); err == nil {
		t.Error("a feed signed with another key was accepted")
	}
	u.PublicKey = pub

	// No download for this platform.
	u.Platform = "plan9-amd64"
	if _, err := u.Check(ctx); err == nil || !strings.Contains(err.Error(), "no download") {
		t.Errorf("missing platform: %v", err)
	}
}

func TestManaged(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux paths")
	}
	if !managed("/usr/bin/app") {
		t.Error("/usr/bin is a package manager's")
	}
	if managed(filepath.Join(t.TempDir(), "app")) {
		t.Error("a user's own directory isn't managed")
	}
	u := &Updater{FeedURL: "http://x/feed.json", PublicKey: "k", Executable: "/usr/bin/app"}
	if _, err := u.Check(context.Background()); err != ErrManaged {
		t.Errorf("Check on /usr/bin: %v, want ErrManaged", err)
	}
}

// TestInstall_Linux replaces a program and restarts it.
func TestInstall_Linux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Linux install")
	}
	dir := t.TempDir()
	exe := filepath.Join(dir, "app")
	marker := filepath.Join(dir, "ran")
	os.WriteFile(exe, []byte("#!/bin/sh\necho old > "+marker+"\n"), 0o755)
	newProg := filepath.Join(t.TempDir(), "download")
	os.WriteFile(newProg, []byte("#!/bin/sh\necho \"new $"+WaitEnv+"\" > "+marker+"\n"), 0o644)

	exited := -1
	exit = func(code int) { exited = code }
	defer func() { exit = os.Exit }()
	if err := install(exe, newProg, "binary"); err != nil {
		t.Fatal(err)
	}
	if exited != 0 {
		t.Error("install didn't exit after restarting")
	}
	if b, _ := os.ReadFile(exe + ".old"); !strings.Contains(string(b), "echo old") {
		t.Error("the old program wasn't kept as .old")
	}
	var ran string
	for deadline := time.Now().Add(5 * time.Second); time.Now().Before(deadline); time.Sleep(20 * time.Millisecond) {
		if b, err := os.ReadFile(marker); err == nil && len(b) > 0 {
			ran = strings.TrimSpace(string(b))
			break
		}
	}
	if !strings.HasPrefix(ran, "new ") || strings.TrimPrefix(ran, "new ") == "" {
		t.Errorf("restarted program wrote %q, want the new one with the old pid", ran)
	}
	if info, _ := os.Stat(exe); info.Mode()&0o100 == 0 {
		t.Error("the new program isn't executable")
	}
}

func TestResolve(t *testing.T) {
	if got, _ := resolve("https://example.com/up/feed.json", "App-1.3.exe"); got != "https://example.com/up/App-1.3.exe" {
		t.Errorf("relative: %s", got)
	}
	if got, _ := resolve("https://example.com/up/feed.json", "https://cdn.example.com/a"); got != "https://cdn.example.com/a" {
		t.Errorf("absolute: %s", got)
	}
}
