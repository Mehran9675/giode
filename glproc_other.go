//go:build !windows && !darwin && !linux && !android && !freebsd && !openbsd && !netbsd

package giode

import "unsafe"

func glGetProcAddress(name string) unsafe.Pointer { return nil }

func glReadPixels(fb uint32, w, h int, out []byte) bool { return false }
