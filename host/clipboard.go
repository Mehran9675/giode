package host

import (
	"io"
	"strings"

	"gioui.org/io/input"
	"gioui.org/io/transfer"
)

// ServeClipboard connects Gio's input router to the host's clipboard
// (GLFW's, the OS's): text components copy (editors, SelectionArea,
// kit.WriteClipboard) goes to write, and their paste requests are
// answered with read's text. Call it every frame, after router.Frame.
// giode.App's window does this itself.
func ServeClipboard(router *input.Router, read func() string, write func(text string)) {
	if _, content, ok := router.WriteClipboard(); ok && write != nil {
		write(string(content))
	}
	if router.ClipboardRequested() && read != nil {
		text := read()
		router.Queue(transfer.DataEvent{
			Type: "application/text",
			Open: func() io.ReadCloser { return io.NopCloser(strings.NewReader(text)) },
		})
	}
}
