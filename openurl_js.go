//go:build js

package giode

import (
	"context"
	"syscall/js"
)

func (a *App) openURL(ctx context.Context, u string) error {
	js.Global().Call("open", u, "_blank", "noopener")
	return nil
}
