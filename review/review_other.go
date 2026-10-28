//go:build !android && !darwin

package review

import (
	"context"

	"github.com/mehran9675/giode"
)

func request(ctx context.Context, a *giode.App) error { return giode.ErrNotSupported }
