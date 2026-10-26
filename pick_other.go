//go:build !android && !ios

package giode

import "context"

func (a *App) pickPhoto(ctx context.Context, opts PhotoOptions) (PickedFile, error) {
	return a.photoFromFile(ctx, opts)
}
