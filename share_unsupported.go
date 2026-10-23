//go:build ios || js

package giode

import "context"

func (t ShareText) share(ctx context.Context, a *App) error { return ErrNotSupported }

func (f ShareFile) share(ctx context.Context, a *App) error { return ErrNotSupported }

func (a *App) watchLifecycle() {}

const lifecycleFromWindowState = false
