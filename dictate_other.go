//go:build !android && !ios

package giode

import "context"

func (a *App) dictate(ctx context.Context, lang string) (<-chan Transcript, error) {
	return nil, ErrNotSupported
}
