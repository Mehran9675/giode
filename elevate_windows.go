package giode

import (
	"errors"
	"strings"

	"golang.org/x/sys/windows"
)

func isElevated() bool {
	return windows.GetCurrentProcessToken().IsElevated()
}

func startElevated(a *App, exe string, args []string) error {
	quoted := make([]string, len(args))
	for i, s := range args {
		quoted[i] = windows.EscapeArg(s)
	}
	err := windows.ShellExecute(windows.Handle(a.NativeHandle()), windows.StringToUTF16Ptr("runas"),
		windows.StringToUTF16Ptr(exe), windows.StringToUTF16Ptr(strings.Join(quoted, " ")), nil, windows.SW_SHOWNORMAL)
	if errors.Is(err, windows.ERROR_CANCELLED) {
		return ErrDeclined
	}
	return err
}
