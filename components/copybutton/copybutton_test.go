package copybutton

import (
	"image"
	"testing"
	"time"

	"gioui.org/io/input"
	"gioui.org/layout"
	"gioui.org/op"
)

func TestCopyButton_CopiesAndConfirms(t *testing.T) {
	var r input.Router
	now := time.Unix(0, 0)
	var copied []string
	b := New(func() string { return "secret" }).OnCopy(func(s string) { copied = append(copied, s) })
	frame := func() {
		var ops op.Ops
		gtx := layout.Context{Ops: &ops, Now: now, Constraints: layout.Constraints{Max: image.Pt(200, 50)}, Source: r.Source()}
		b.Layout(gtx)
		r.Frame(&ops)
	}
	frame()
	b.Click()
	frame()
	if _, data, ok := r.WriteClipboard(); !ok || string(data) != "secret" {
		t.Fatalf("clipboard = %q, %v", data, ok)
	}
	if len(copied) != 1 {
		t.Errorf("OnCopy calls = %v", copied)
	}
	now = now.Add(100 * time.Millisecond)
	frame()
	if !b.Copied() {
		t.Error("confirmation not showing right after copying")
	}
	now = now.Add(2 * time.Second)
	frame()
	if b.Copied() {
		t.Error("confirmation still showing after the feedback time")
	}
}
