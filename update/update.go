// Package update keeps a desktop app up to date from a feed you host:
// it checks for a newer version, downloads it, verifies it, replaces
// the app and restarts it.
//
//	u := update.New() // feed, key and version from giode.toml, via the giode CLI
//	rel, err := u.Check(ctx)
//	if err == nil && rel != nil {
//		dl, err := u.Download(ctx, rel, nil)
//		if err == nil {
//			err = dl.Install() // restarts into the new version
//		}
//	}
//
// The feed is a JSON file listing the latest version and a download per
// platform, signed with an Ed25519 key: `giode update-key` makes the key
// pair (the public half goes into the app), and `giode release` writes
// the signed feed and the downloads from a build. The app trusts only
// a feed signed with its key, and only downloads whose SHA-256 the feed
// lists, so the server needn't be trusted.
//
// Apps installed by a package manager or store - a .deb or .rpm,
// Flatpak, MSIX, the App Store - are updated by it instead: Check
// returns ErrManaged for them.
package update

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// Set by the giode CLI from giode.toml's [update] and app.version.
var (
	defaultFeedURL   string
	defaultPublicKey string
	defaultVersion   string
)

// ErrManaged is returned for an app a package manager or store
// installed and updates.
var ErrManaged = errors.New("update: the app is installed by a package manager or store, which updates it")

// Updater checks a feed and installs its releases.
type Updater struct {
	// FeedURL is the feed's address; its signature is at FeedURL+".sig".
	FeedURL string
	// PublicKey is the base64 Ed25519 public key the feed is signed
	// with.
	PublicKey string
	// Version is the running version, major.minor.patch.
	Version string
	// Client fetches the feed and downloads; http.DefaultClient when
	// nil.
	Client *http.Client
	// Platform is the feed's key for this build ("windows-amd64",
	// "darwin-universal", "linux-amd64"); detected when empty.
	Platform string
	// Executable is the program (or .app bundle, or AppImage) to
	// replace; detected when empty.
	Executable string
}

// New returns an Updater configured by the giode CLI from giode.toml.
func New() *Updater {
	return &Updater{FeedURL: defaultFeedURL, PublicKey: defaultPublicKey, Version: defaultVersion}
}

// Feed is the update feed.
type Feed struct {
	// Version is the latest version; Notes describe it; Date is when
	// it was released (RFC 3339).
	Version string `json:"version"`
	Notes   string `json:"notes,omitempty"`
	Date    string `json:"date,omitempty"`
	// Assets are the downloads by platform.
	Assets map[string]Asset `json:"assets"`
}

// Asset is a download for one platform.
type Asset struct {
	// URL is where it is, relative to the feed or absolute.
	URL string `json:"url"`
	// SHA256 is the hex digest of the file; Size its length.
	SHA256 string `json:"sha256"`
	Size   int64  `json:"size"`
	// Kind says how it installs: "exe" (a Windows program replacing
	// the running one), "installer" (a Windows setup program run
	// silently), "app.zip" (a zipped macOS app bundle), "binary" (a
	// Linux program) or "appimage".
	Kind string `json:"kind"`
}

// Release is an available update.
type Release struct {
	Version, Notes, Date string
	Asset                Asset
	// url is the asset's absolute address.
	url string
}

// Check fetches the feed and returns the release newer than Version
// for this platform, or nil when the app is up to date.
func (u *Updater) Check(ctx context.Context) (*Release, error) {
	if u.FeedURL == "" || u.PublicKey == "" {
		return nil, errors.New("update: no feed URL or public key (set [update] in giode.toml, or Updater's fields)")
	}
	if managed(u.executable()) {
		return nil, ErrManaged
	}
	feedBytes, err := u.get(ctx, u.FeedURL, 1<<20)
	if err != nil {
		return nil, err
	}
	sig, err := u.get(ctx, u.FeedURL+".sig", 4096)
	if err != nil {
		return nil, err
	}
	if err := Verify(u.PublicKey, feedBytes, sig); err != nil {
		return nil, err
	}
	var f Feed
	if err := json.Unmarshal(feedBytes, &f); err != nil {
		return nil, fmt.Errorf("update: bad feed: %w", err)
	}
	if !Newer(f.Version, u.Version) {
		return nil, nil
	}
	a, ok := f.Assets[u.platform()]
	if !ok {
		return nil, fmt.Errorf("update: version %s has no download for %s", f.Version, u.platform())
	}
	abs, err := resolve(u.FeedURL, a.URL)
	if err != nil {
		return nil, err
	}
	return &Release{Version: f.Version, Notes: f.Notes, Date: f.Date, Asset: a, url: abs}, nil
}

// Download is a downloaded, verified release.
type Download struct {
	u       *Updater
	Release *Release
	// Path is the downloaded file.
	Path string
}

// Download fetches the release into a temporary file, checking its
// SHA-256. progress, when set, is called as bytes arrive.
func (u *Updater) Download(ctx context.Context, r *Release, progress func(done, total int64)) (*Download, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := u.client().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update: %s: %s", r.url, resp.Status)
	}
	f, err := os.CreateTemp("", "giode-update-*"+filepath.Ext(strings.SplitN(r.url, "?", 2)[0]))
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	var done int64
	buf := make([]byte, 64<<10)
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			done += int64(n)
			if r.Asset.Size > 0 && done > r.Asset.Size {
				f.Close()
				os.Remove(f.Name())
				return nil, errors.New("update: the download is larger than the feed says")
			}
			h.Write(buf[:n])
			if _, err := f.Write(buf[:n]); err != nil {
				f.Close()
				os.Remove(f.Name())
				return nil, err
			}
			if progress != nil {
				progress(done, r.Asset.Size)
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			f.Close()
			os.Remove(f.Name())
			return nil, rerr
		}
	}
	f.Close()
	if got := hex.EncodeToString(h.Sum(nil)); !strings.EqualFold(got, r.Asset.SHA256) {
		os.Remove(f.Name())
		return nil, fmt.Errorf("update: the download's SHA-256 is %s, the feed says %s", got, r.Asset.SHA256)
	}
	return &Download{u: u, Release: r, Path: f.Name()}, nil
}

// Install replaces the app with the download and restarts it; on
// success it doesn't return. Close windows and save state first.
func (d *Download) Install() error {
	return install(d.u.executable(), d.Path, d.Release.Asset.Kind)
}

// Cleanup removes what an earlier update left behind (the replaced
// program). Call it once at startup.
func Cleanup() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	target := appRoot(exe)
	os.RemoveAll(target + ".old")
}

// Verify checks sig, a base64 Ed25519 signature, over data with the
// base64 public key.
func Verify(publicKey string, data, sig []byte) error {
	pk, err := base64.StdEncoding.DecodeString(strings.TrimSpace(publicKey))
	if err != nil || len(pk) != ed25519.PublicKeySize {
		return errors.New("update: the public key isn't a base64 Ed25519 key")
	}
	s, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(sig)))
	if err != nil || !ed25519.Verify(pk, data, s) {
		return errors.New("update: the feed's signature doesn't match the app's key")
	}
	return nil
}

// Sign signs data with a base64 Ed25519 private key, returning the
// base64 signature the feed's .sig file holds.
func Sign(privateKey string, data []byte) (string, error) {
	k, err := base64.StdEncoding.DecodeString(strings.TrimSpace(privateKey))
	if err != nil || len(k) != ed25519.PrivateKeySize {
		return "", errors.New("update: the private key isn't a base64 Ed25519 key")
	}
	return base64.StdEncoding.EncodeToString(ed25519.Sign(ed25519.PrivateKey(k), data)), nil
}

// Newer reports whether version a is newer than b (major.minor.patch,
// missing parts 0; anything after - or + is ignored).
func Newer(a, b string) bool {
	pa, pb := parts(a), parts(b)
	for i := range pa {
		if pa[i] != pb[i] {
			return pa[i] > pb[i]
		}
	}
	return false
}

func parts(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	if i := strings.IndexAny(v, "-+"); i >= 0 {
		v = v[:i]
	}
	var out [3]int
	for i, p := range strings.SplitN(v, ".", 3) {
		out[i], _ = strconv.Atoi(p)
	}
	return out
}

func (u *Updater) client() *http.Client {
	if u.Client != nil {
		return u.Client
	}
	return http.DefaultClient
}

func (u *Updater) get(ctx context.Context, url string, limit int64) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := u.client().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update: %s: %s", url, resp.Status)
	}
	return io.ReadAll(io.LimitReader(resp.Body, limit))
}

func (u *Updater) platform() string {
	if u.Platform != "" {
		return u.Platform
	}
	return Platform()
}

// Platform returns this build's feed key: GOOS-GOARCH, or
// darwin-universal on macOS.
func Platform() string {
	if runtime.GOOS == "darwin" {
		return "darwin-universal"
	}
	return runtime.GOOS + "-" + runtime.GOARCH
}

func (u *Updater) executable() string {
	if u.Executable != "" {
		return u.Executable
	}
	if p := os.Getenv("APPIMAGE"); p != "" && runtime.GOOS == "linux" {
		return p
	}
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	if r, err := filepath.EvalSymlinks(exe); err == nil {
		exe = r
	}
	return appRoot(exe)
}

// appRoot is what an update replaces: the .app bundle holding a macOS
// program, else the program.
func appRoot(exe string) string {
	if runtime.GOOS == "darwin" {
		if i := strings.Index(exe, ".app/Contents/MacOS/"); i >= 0 {
			return exe[:i+len(".app")]
		}
	}
	return exe
}

// managed reports whether the app at path belongs to a package
// manager or store.
func managed(path string) bool {
	switch {
	case path == "":
		return true
	case os.Getenv("FLATPAK_ID") != "" || os.Getenv("SNAP") != "":
		return true
	case runtime.GOOS == "linux" && (strings.HasPrefix(path, "/usr/") || strings.HasPrefix(path, "/opt/") && !writable(path)):
		return true
	case runtime.GOOS == "windows" && strings.Contains(strings.ToLower(path), `\windowsapps\`):
		return true
	case runtime.GOOS == "darwin" && fileExists(filepath.Join(path, "Contents", "_MASReceipt")):
		return true
	}
	return false
}

func writable(path string) bool {
	f, err := os.CreateTemp(filepath.Dir(path), ".giode-write-test-*")
	if err != nil {
		return false
	}
	f.Close()
	os.Remove(f.Name())
	return true
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func resolve(base, ref string) (string, error) {
	if strings.Contains(ref, "://") {
		return ref, nil
	}
	i := strings.LastIndex(base, "/")
	if i < 0 {
		return "", fmt.Errorf("update: can't resolve %q against %q", ref, base)
	}
	return base[:i+1] + ref, nil
}
