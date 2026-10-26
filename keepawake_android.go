package giode

import (
	"context"

	"gioui.org/app"

	"github.com/mehran9675/giode/internal/mobile"
)

func (a *App) haptic(k HapticKind) {
	a.onAndroidMain(context.Background(), func(view uintptr) error {
		return mobile.Haptic(app.JavaVM(), app.AppContext(), view, androidHaptic[k])
	})
}

func (a *App) keepScreenOn(ctx context.Context, on bool) error {
	return a.onAndroidMain(ctx, func(view uintptr) error {
		return mobile.KeepScreenOn(app.JavaVM(), app.AppContext(), view, on)
	})
}

func (a *App) keepAlive(reason string) (func(), error) {
	return func() {}, ErrNotSupported
}
