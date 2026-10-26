package segmented

import (
	"image"
	"testing"
	"time"

	"gioui.org/f32"
	"gioui.org/io/input"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/fonts"
	"github.com/mehran9675/giode/styles"
)

func TestSegmented(t *testing.T) {
	fonts.UseGofont()
	var picked []int
	c := New("Chat", "Research", "Code").Styles(styles.Styles{Width: -1}).OnChange(func(i int) { picked = append(picked, i) })
	var r input.Router
	now := time.Unix(0, 0)
	var dims layout.Dimensions
	frame := func() {
		now = now.Add(50 * time.Millisecond)
		var ops op.Ops
		dims = c.Layout(layout.Context{Ops: &ops, Now: now, Source: r.Source(), Constraints: layout.Constraints{Max: image.Pt(300, 100)}})
		r.Frame(&ops)
	}
	frame()
	if dims.Size.X != 300 {
		t.Fatalf("width %d, want the full 300", dims.Size.X)
	}
	p := f32.Pt(250, float32(dims.Size.Y/2)) // the third segment
	r.Queue(pointer.Event{Kind: pointer.Press, Source: pointer.Touch, Position: p, Buttons: pointer.ButtonPrimary},
		pointer.Event{Kind: pointer.Release, Source: pointer.Touch, Position: p})
	for i := 0; i < 8; i++ {
		frame()
	}
	if c.Selected() != 2 || len(picked) != 1 || picked[0] != 2 {
		t.Errorf("selected %d, picked %v; want 2", c.Selected(), picked)
	}
	if c.pos != 2 {
		t.Errorf("highlight at %v, want 2", c.pos)
	}
}
