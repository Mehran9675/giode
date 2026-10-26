//go:build js

package giode

import "context"

func (a *App) haptic(HapticKind) {}

func (a *App) keepScreenOn(ctx context.Context, on bool) error { return ErrNotSupported }

func (a *App) keepAlive(reason string) (func(), error) { return func() {}, ErrNotSupported }
