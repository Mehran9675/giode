package host

import (
	"image"
	"io"
	"strings"
	"testing"

	"gioui.org/io/clipboard"
	"gioui.org/io/event"
	"gioui.org/io/input"
	"gioui.org/io/pointer"
	"gioui.org/io/transfer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
)

// TestServeClipboard copies and pastes through a fake host clipboard.
func TestServeClipboard(t *testing.T) {
	var r input.Router
	board := "from the host"
	var written string
	tag := new(int)
	frame := func(f func(gtx layout.Context)) {
		var ops op.Ops
		gtx := layout.Context{Ops: &ops, Source: r.Source(), Constraints: layout.Exact(image.Pt(10, 10))}
		ps := pointer.PassOp{}.Push(gtx.Ops)
		cl := clip.Rect{Max: image.Pt(10, 10)}.Push(gtx.Ops)
		event.Op(gtx.Ops, tag)
		cl.Pop()
		ps.Pop()
		f(gtx)
		r.Frame(&ops)
		ServeClipboard(&r, func() string { return board }, func(s string) { written = s })
	}
	frame(func(gtx layout.Context) {
		gtx.Execute(clipboard.WriteCmd{Type: "application/text", Data: io.NopCloser(strings.NewReader("copied"))})
		gtx.Execute(clipboard.ReadCmd{Tag: tag})
	})
	if written != "copied" {
		t.Errorf("host clipboard got %q", written)
	}
	got := ""
	frame(func(gtx layout.Context) {
		for {
			ev, ok := gtx.Event(transfer.TargetFilter{Target: tag, Type: "application/text"})
			if !ok {
				break
			}
			if e, ok := ev.(transfer.DataEvent); ok {
				b, _ := io.ReadAll(e.Open())
				got = string(b)
			}
		}
	})
	if got != board {
		t.Errorf("paste got %q, want %q", got, board)
	}
}
