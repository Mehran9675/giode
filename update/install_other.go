//go:build !(linux && !android) && !(darwin && !ios) && !windows

package update

import "fmt"

func install(target, download, kind string) error {
	return fmt.Errorf("update: apps on this platform are updated by its store")
}
