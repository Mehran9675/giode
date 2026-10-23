//go:build !android

package mobile

// WatchLifecycle reports ErrUnsupported outside Android.
func WatchLifecycle(vm, ctx uintptr) error { return ErrUnsupported }

// Publish reports ErrUnsupported outside Android.
func Publish(vm, ctx uintptr, path, mime string) (string, error) { return "", ErrUnsupported }

// ShareText reports ErrUnsupported outside Android.
func ShareText(vm, ctx, view uintptr, text, subject, title string) error { return ErrUnsupported }

// ShareURI reports ErrUnsupported outside Android.
func ShareURI(vm, ctx, view uintptr, uri, mime, title string) error { return ErrUnsupported }
