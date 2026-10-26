//go:build !linux && !darwin && !windows

package giode

func deviceDetails() (ram uint64, version, model string) { return 0, "", "" }
