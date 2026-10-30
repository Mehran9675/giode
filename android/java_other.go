//go:build !android

package android

import (
	"errors"
	"unsafe"
)

// Available reports whether Java can be called: true on Android.
func Available() bool { return false }

// CallVoid calls a static void method.
func CallVoid(class, method, sig string, args ...any) error { return errors.ErrUnsupported }

// CallString calls a static method returning a String; null is "".
func CallString(class, method, sig string, args ...any) (string, error) {
	return "", errors.ErrUnsupported
}

// CallBool calls a static method returning a boolean.
func CallBool(class, method, sig string, args ...any) (bool, error) {
	return false, errors.ErrUnsupported
}

// CallInt calls a static method returning an int.
func CallInt(class, method, sig string, args ...any) (int, error) {
	return 0, errors.ErrUnsupported
}

// CallBytes calls a static method returning a byte[]; null is nil.
func CallBytes(class, method, sig string, args ...any) ([]byte, error) {
	return nil, errors.ErrUnsupported
}

// GoString converts a Java string passed to a native method; null is
// "".
func GoString(env unsafe.Pointer, s uintptr) string { return "" }

// GoBytes converts a Java byte[] passed to a native method; null is nil.
func GoBytes(env unsafe.Pointer, b uintptr) []byte { return nil }
