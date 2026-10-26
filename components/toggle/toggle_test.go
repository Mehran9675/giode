package toggle

import (
	"image"
	"testing"
	"time"

	"gioui.org/f32"
	"gioui.org/io/input"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
)

func TestSwitch(t *testing.T) {
	var changes []bool
	s := New("").OnChange(func(on bool) { changes = append(changes, on) })
	var r input.Router
	now := time.Unix(0, 0)
	var dims layout.Dimensions
	frame := func() {
		now = now.Add(50 * time.Millisecond)
		var ops op.Ops
		dims = s.Layout(layout.Context{Ops: &ops, Now: now, Source: r.Source(), Constraints: layout.Constraints{Max: image.Pt(400, 100)}})
		r.Frame(&ops)
	}
	frame()
	if dims.Size.X < 40 || dims.Size.Y < 20 {
		t.Fatalf("switch size %v", dims.Size)
	}
	p := f32.Pt(10, 10)
	r.Queue(pointer.Event{Kind: pointer.Press, Source: pointer.Mouse, Position: p, Buttons: pointer.ButtonPrimary},
		pointer.Event{Kind: pointer.Release, Source: pointer.Mouse, Position: p})
	for i := 0; i < 6; i++ {
		frame()
	}
	if !s.Checked() || len(changes) != 1 || !changes[0] {
		t.Fatalf("after a tap: checked=%v changes=%v", s.Checked(), changes)
	}
	if s.pos != 1 {
		t.Errorf("knob at %v after the animation, want 1", s.pos)
	}
	s.Disabled(true).Toggle()
	if !s.Checked() {
		t.Error("a disabled switch toggled")
	}
}
