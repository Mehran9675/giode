package elements

import (
	"image"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

func TestMinTouchTarget(t *testing.T) {
	defer SetTouchInput(TouchInput())
	small := func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: image.Pt(20, 10), Baseline: 2}
	}
	var ops op.Ops
	gtx := layout.Context{Ops: &ops, Metric: unit.Metric{PxPerDp: 2}, Constraints: layout.Constraints{Max: image.Pt(500, 500)}}

	SetTouchInput(false)
	if got := MinTouchTarget(gtx, small).Size; got != image.Pt(20, 10) {
		t.Errorf("mouse: size %v, want the content's 20x10", got)
	}

	SetTouchInput(true)
	dims := MinTouchTarget(gtx, small)
	if dims.Size != image.Pt(96, 96) {
		t.Errorf("touch: size %v, want 48dp = 96px square", dims.Size)
	}
	// Centered: 43px below the content, which kept its 2px baseline.
	if dims.Baseline != 45 {
		t.Errorf("touch: baseline %d, want 45", dims.Baseline)
	}

	// Never beyond the constraints, and larger content is unchanged.
	gtx.Constraints.Max = image.Pt(60, 500)
	if got := MinTouchTarget(gtx, small).Size; got != image.Pt(60, 96) {
		t.Errorf("narrow: size %v, want 60x96", got)
	}
	big := func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{Size: image.Pt(200, 120)} }
	gtx.Constraints.Max = image.Pt(500, 500)
	if got := MinTouchTarget(gtx, big).Size; got != image.Pt(200, 120) {
		t.Errorf("big: size %v, want 200x120", got)
	}
}

func TestSetTouchInput_ReportsChange(t *testing.T) {
	defer SetTouchInput(TouchInput())
	SetTouchInput(false)
	if !SetTouchInput(true) || SetTouchInput(true) {
		t.Error("SetTouchInput change reporting wrong")
	}
}
