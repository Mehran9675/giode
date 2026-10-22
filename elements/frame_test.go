package elements

import (
	"testing"

	"gioui.org/unit"
)

func TestAfterLayoutRunsOnceInOrder(t *testing.T) {
	var got []int
	AfterLayout(func() { got = append(got, 1) })
	AfterLayout(func() { got = append(got, 2) })
	if !RunAfterLayout() {
		t.Fatal("RunAfterLayout reported nothing ran")
	}
	if RunAfterLayout() {
		t.Error("queue not cleared")
	}
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Errorf("ran %v, want [1 2]", got)
	}
}

func TestDisplayScale(t *testing.T) {
	defer SetMetric(unit.Metric{})
	SetMetric(unit.Metric{PxPerDp: 1.5})
	if DisplayScale() != 1.5 || DpToPx(10) != 15 || PxToDp(30) != 20 {
		t.Errorf("scale %v, DpToPx(10) %d, PxToDp(30) %v", DisplayScale(), DpToPx(10), PxToDp(30))
	}
	SetMetric(unit.Metric{})
	if DisplayScale() != 1 {
		t.Errorf("zero metric scale = %v, want 1", DisplayScale())
	}
}
