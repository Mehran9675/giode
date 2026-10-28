//go:build !giode_gio2

package giode

import (
	"context"

	"gioui.org/io/event"
)

func injectEvent(ctx context.Context, a *App, e event.Event) error {
	return errInjectUnpatched
}
