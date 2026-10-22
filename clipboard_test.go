package giode

import (
	"image"
	"io"
	"strings"
	"testing"

	"gioui.org/app"
	"gioui.org/io/input"
	"gioui.org/io/transfer"
	"gioui.org/layout"
	"gioui.org/op"
)

func TestClipboard_WriteAndRead(t *testing.T) {
	var r input.Router
	a := &App{window: new(app.Window)}
	frame := func() {
		var ops op.Ops
		gtx := layout.Context{Ops: &ops, Constraints: layout.Constraints{Max: image.Pt(100, 100)}, Source: r.Source()}
		a.layoutClipboard(gtx)
		r.Frame(&ops)
	}

	a.SetClipboard("hello")
	frame()
	if mime, data, ok := r.WriteClipboard(); !ok || mime != "application/text" || string(data) != "hello" {
		t.Errorf("clipboard write = %q %q %v", mime, data, ok)
	}

	var got []string
	a.ReadClipboard(func(s string) { got = append(got, s) })
	a.ReadClipboard(func(s string) { got = append(got, s+"!") })
	frame()
	if !r.ClipboardRequested() {
		t.Fatal("no clipboard read requested")
	}
	r.Queue(transfer.DataEvent{Type: "application/text", Open: func() io.ReadCloser {
		return io.NopCloser(strings.NewReader("pasted"))
	}})
	frame()
	if len(got) != 2 || got[0] != "pasted" || got[1] != "pasted!" {
		t.Errorf("readers got %v", got)
	}
}
