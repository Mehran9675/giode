//go:build (linux && !android) || freebsd || openbsd || netbsd || dragonfly

package host

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/godbus/dbus/v5"
)

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

// KeepScreenOn keeps the screen from dimming and locking while on is
// set: the freedesktop screen saver inhibition.
func KeepScreenOn(on bool) error {
	return keepScreenOn(context.Background(), on)
}

func keepScreenOn(ctx context.Context, on bool) error {
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

// KeepAwake keeps the computer from sleeping while idle until release
// is called: a systemd-logind inhibitor lock, listed with reason.
func KeepAwake(reason string) (release func(), err error) {
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
