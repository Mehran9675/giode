//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package giode

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/godbus/dbus/v5"
)

func (a *App) haptic(HapticKind) {}

// screenCookie is the freedesktop screen saver inhibition held by
// KeepScreenOn.
var screen struct {
	sync.Mutex
	cookie uint32
	held   bool
}

func appName() string {
	exe, err := os.Executable()
	if err != nil {
		return "giode"
	}
	return filepath.Base(exe)
}

func (a *App) keepScreenOn(ctx context.Context, on bool) error {
	screen.Lock()
	defer screen.Unlock()
	if on == screen.held {
		return nil
	}
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	obj := conn.Object("org.freedesktop.ScreenSaver", "/org/freedesktop/ScreenSaver")
	if on {
		if err := obj.CallWithContext(ctx, "org.freedesktop.ScreenSaver.Inhibit", 0, appName(), "Keeping the screen on").Store(&screen.cookie); err != nil {
			return err
		}
	} else if err := obj.CallWithContext(ctx, "org.freedesktop.ScreenSaver.UnInhibit", 0, screen.cookie).Err; err != nil {
		return err
	}
	screen.held = on
	return nil
}

// keepAlive takes a systemd-logind inhibitor lock against idle sleep;
// closing its file descriptor releases it.
func (a *App) keepAlive(reason string) (func(), error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return func() {}, err
	}
	var fd dbus.UnixFD
	err = conn.Object("org.freedesktop.login1", "/org/freedesktop/login1").
		Call("org.freedesktop.login1.Manager.Inhibit", 0, "sleep:idle", appName(), reason, "block").Store(&fd)
	if err != nil {
		return func() {}, err
	}
	f := os.NewFile(uintptr(fd), "inhibit")
	var once sync.Once
	return func() { once.Do(func() { f.Close() }) }, nil
}
