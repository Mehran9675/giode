package tooltip

import (
	"image"
	"testing"

	"gioui.org/f32"
	"gioui.org/io/input"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
)

func TestTooltip_HoveredFollowsPointer(t *testing.T) {
	elements.SetWindowSize(image.Pt(800, 600))
	defer elements.SetWindowSize(image.Point{})

	var ops op.Ops
	var router input.Router
	gtx := layout.Context{Ops: &ops, Source: router.Source(), Constraints: layout.Exact(image.Pt(100, 40))}

	trigger := elements.Box(styles.Styles{Width: 100, Height: 40})
	tip := New(trigger, "hint")

	moveTo := func(p f32.Point) {
		router.Queue(pointer.Event{Kind: pointer.Move, Source: pointer.Mouse, Position: p})
	}
	frame := func() {
		ops.Reset()
		tip.Layout(gtx)
		router.Frame(&ops)
	}

	frame() // register the hover region

	moveTo(f32.Pt(10, 10))
	frame()
	frame()
	if !tip.Hovered() {
		t.Fatalf("pointer inside trigger: Hovered() = false, want true")
	}

	moveTo(f32.Pt(500, 500))
	frame()
	frame()
	if tip.Hovered() {
		t.Fatalf("pointer outside trigger: Hovered() = true, want false")
	}
}

func TestTooltip_LayoutReturnsTriggerDimensions(t *testing.T) {
	elements.SetWindowSize(image.Pt(800, 600))
	defer elements.SetWindowSize(image.Point{})

	var ops op.Ops
	var router input.Router
	gtx := layout.Context{Ops: &ops, Source: router.Source(), Constraints: layout.Exact(image.Pt(100, 40))}

	trigger := elements.Box(styles.Styles{Width: 100, Height: 40})
	tip := New(trigger, "hint")

	dims := tip.Layout(gtx)
	if dims.Size != (image.Point{X: 100, Y: 40}) {
		t.Errorf("Layout dims = %v, want the trigger's own 100x40", dims.Size)
	}
}
