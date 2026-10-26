//go:build linux && !android

package giode

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

// systemDark reads the XDG desktop portal's appearance setting:
// color-scheme 1 means dark.
func systemDark() (bool, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return false, errors.Join(ErrNotSupported, err)
	}
	var v dbus.Variant
	err = conn.Object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop").
		Call("org.freedesktop.portal.Settings.Read", 0, "org.freedesktop.appearance", "color-scheme").Store(&v)
	if err != nil {
		return false, errors.Join(ErrNotSupported, err)
	}
	// Read wraps the value in a second variant.
	if inner, ok := v.Value().(dbus.Variant); ok {
		v = inner
	}
	scheme, ok := v.Value().(uint32)
	if !ok {
		return false, ErrNotSupported
	}
	return scheme == 1, nil
}
