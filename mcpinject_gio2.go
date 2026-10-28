//go:build giode_gio2

package giode

import (
	"context"

	"gioui.org/io/event"
)

// injectEvent delivers e to the window's input through giode's patched
// Gio. It must not run on the UI goroutine.
func injectEvent(ctx context.Context, a *App, e event.Event) error {
	done := make(chan struct{})
	go func() {
		a.window.GiodeInject(e)
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
