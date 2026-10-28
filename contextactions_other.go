//go:build !windows && !(linux && !android) && !(darwin && !ios)

package giode

import "errors"

func installContextActions(string, string, []ContextAction) error {
	return errors.New("giode: context actions are a desktop feature")
}

func uninstallContextActions(string) error { return nil }
