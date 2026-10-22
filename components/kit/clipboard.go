package kit

import (
	"io"
	"strings"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
)

// WriteClipboard copies text to the system clipboard from inside a
// Layout. Outside layout, use App.SetClipboard.
func WriteClipboard(gtx layout.Context, text string) {
	gtx.Execute(clipboard.WriteCmd{Type: "application/text", Data: io.NopCloser(strings.NewReader(text))})
}
