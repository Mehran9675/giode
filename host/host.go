// Package host has the platform services of giode.App that don't need
// its window, for apps that host Giode themselves - in a GLFW window,
// as an embedded view (see docs/low-level.md): keeping the screen on,
// the single-instance lock, and the clipboard bridge between Gio's
// input router and the host's clipboard. App uses them too. The
// package doesn't import gioui.org/app.
package host

import "errors"

// ErrNotSupported is returned by features the current platform lacks.
// giode.ErrNotSupported is the same error.
var ErrNotSupported = errors.New("giode: not supported on this platform")
