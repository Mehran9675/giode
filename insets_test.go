package giode

import (
	"image"
	"testing"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
)

func TestFrameContext_SafeAreaAndEdgeToEdge(t *testing.T) {
	e := app.FrameEvent{
		Size:   image.Pt(1080, 2400),
		Metric: unit.Metric{PxPerDp: 3},
		Insets: app.Insets{Top: 24, Bottom: 48},
	}
	var ops op.Ops
	a := &App{}
	if got := a.frameContext(&ops, e).Constraints.Max; got != image.Pt(1080, 2400-72*3) {
		t.Errorf("inside the safe area: %v, want the height minus 72dp", got)
	}
	if in := SafeArea(); in.Top != 24 || in.Bottom != 48 {
		t.Errorf("SafeArea = %+v, want top 24 and bottom 48", in)
	}
	a.cfg.EdgeToEdge = true
	if got := a.frameContext(&ops, e).Constraints.Max; got != e.Size {
		t.Errorf("edge to edge: %v, want the whole window", got)
	}
}

func TestDevice(t *testing.T) {
	d := Device()
	if d.OS == "" || d.Class != DeviceDesktop {
		t.Errorf("Device() = %+v", d)
	}
	if d.TotalRAM == 0 {
		t.Error("TotalRAM unknown on Linux")
	}
}
