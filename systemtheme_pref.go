//go:build !linux || android

package giode

import (
	"errors"

	"gioui.org/x/pref/theme"
)

func systemDark() (bool, error) {
	dark, err := theme.IsDarkMode()
	if err != nil {
		return false, errors.Join(ErrNotSupported, err)
	}
	return dark, nil
}
