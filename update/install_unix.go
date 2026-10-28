//go:build (linux && !android) || (darwin && !ios)

package update

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func install(target, download, kind string) error {
	switch {
	case runtime.GOOS == "linux" && (kind == "binary" || kind == "appimage"):
		if err := replaceFile(target, download); err != nil {
			return err
		}
		return restart(target)
	case runtime.GOOS == "darwin" && kind == "app.zip":
		if !strings.HasSuffix(target, ".app") {
			return fmt.Errorf("update: %s isn't in an app bundle", target)
		}
		tmp, err := sameDirTemp(target)
		if err != nil {
			return fmt.Errorf("update: %w", err)
		}
		defer os.RemoveAll(tmp)
		app, err := unzipApp(download, tmp)
		if err != nil {
			return err
		}
		os.RemoveAll(target + ".old")
		if err := os.Rename(target, target+".old"); err != nil {
			return fmt.Errorf("update: %w", err)
		}
		if err := os.Rename(app, target); err != nil {
			os.Rename(target+".old", target)
			return fmt.Errorf("update: %w", err)
		}
		os.Remove(download)
		exe := filepath.Join(target, "Contents", "MacOS", filepath.Base(os.Args[0]))
		if _, err := os.Stat(exe); err != nil {
			// The program's name changed: take the bundle's only one.
			entries, _ := os.ReadDir(filepath.Join(target, "Contents", "MacOS"))
			if len(entries) != 1 {
				return fmt.Errorf("update: can't find the program in %s", target)
			}
			exe = filepath.Join(target, "Contents", "MacOS", entries[0].Name())
		}
		return restart(exe)
	}
	return fmt.Errorf("update: can't install a %q download on %s", kind, runtime.GOOS)
}

// unzipApp extracts a zipped .app into dir and returns its path,
// keeping permissions and symlinks (frameworks use them).
func unzipApp(zipPath, dir string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", fmt.Errorf("update: %w", err)
	}
	defer r.Close()
	app := ""
	for _, f := range r.File {
		name := filepath.FromSlash(f.Name)
		if strings.Contains(name, "..") || filepath.IsAbs(name) {
			return "", fmt.Errorf("update: bad path %q in the archive", f.Name)
		}
		p := filepath.Join(dir, name)
		if top := strings.SplitN(filepath.ToSlash(name), "/", 2)[0]; strings.HasSuffix(top, ".app") {
			app = filepath.Join(dir, top)
		}
		mode := f.Mode()
		switch {
		case f.FileInfo().IsDir():
			if err := os.MkdirAll(p, 0o755); err != nil {
				return "", err
			}
			continue
		case mode&os.ModeSymlink != 0:
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			link, _ := io.ReadAll(rc)
			rc.Close()
			os.MkdirAll(filepath.Dir(p), 0o755)
			if err := os.Symlink(string(link), p); err != nil {
				return "", err
			}
			continue
		}
		os.MkdirAll(filepath.Dir(p), 0o755)
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		out, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode.Perm()|0o200)
		if err != nil {
			rc.Close()
			return "", err
		}
		_, err = io.Copy(out, rc)
		rc.Close()
		out.Close()
		if err != nil {
			return "", err
		}
	}
	if app == "" {
		return "", fmt.Errorf("update: the archive has no .app")
	}
	return app, nil
}
