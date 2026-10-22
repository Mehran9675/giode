package pagination

import (
	"reflect"
	"testing"

	"gioui.org/f32"
	"gioui.org/io/input"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"image"
)

func TestPagination_WindowedSmallTotal(t *testing.T) {
	p := New(5)
	got := p.windowed()
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("windowed() with total=5 = %v, want %v", got, want)
	}
}

func TestPagination_WindowedLargeTotalNearStart(t *testing.T) {
	p := New(20)
	p.SetCurrent(1)
	got := p.windowed()
	want := []int{1, 2, 0, 20} // hi clamps at current+1=2, no ellipsis needed before, one after
	if !reflect.DeepEqual(got, want) {
		t.Errorf("windowed() at page 1 of 20 = %v, want %v", got, want)
	}
}

func TestPagination_WindowedLargeTotalMiddle(t *testing.T) {
	p := New(20)
	p.SetCurrent(10)
	got := p.windowed()
	want := []int{1, 0, 9, 10, 11, 0, 20}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("windowed() at page 10 of 20 = %v, want %v", got, want)
	}
}

func TestPagination_WindowedLargeTotalNearEnd(t *testing.T) {
	p := New(20)
	p.SetCurrent(20)
	got := p.windowed()
	want := []int{1, 0, 19, 20}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("windowed() at page 20 of 20 = %v, want %v", got, want)
	}
}

func TestPagination_SetCurrentClamps(t *testing.T) {
	p := New(5)
	p.SetCurrent(100)
	if p.Current() != 5 {
		t.Errorf("SetCurrent(100) with total=5: Current() = %d, want 5", p.Current())
	}
	p.SetCurrent(-3)
	if p.Current() != 1 {
		t.Errorf("SetCurrent(-3): Current() = %d, want 1", p.Current())
	}
}

func TestPagination_SetTotalClampsCurrent(t *testing.T) {
	p := New(10)
	p.SetCurrent(8)
	p.SetTotal(5)
	if p.Current() != 5 {
		t.Errorf("SetTotal(5) after Current=8: Current() = %d, want 5", p.Current())
	}
}

func TestPagination_ClickPrevChangesCurrentAndFiresOnChange(t *testing.T) {
	var changedTo int
	p := New(5).OnChange(func(page int) { changedTo = page })
	p.SetCurrent(2)

	var ops op.Ops
	var router input.Router
	gtx := layout.Context{Ops: &ops, Source: router.Source(), Constraints: layout.Exact(image.Pt(400, 60))}

	frame := func() {
		ops.Reset()
		p.Layout(gtx)
		router.Frame(&ops)
	}
	press := func(pt f32.Point) {
		router.Queue(pointer.Event{Kind: pointer.Press, Source: pointer.Mouse, Position: pt, Buttons: pointer.ButtonPrimary})
	}
	release := func(pt f32.Point) {
		router.Queue(pointer.Event{Kind: pointer.Release, Source: pointer.Mouse, Position: pt, Buttons: 0})
	}

	frame() // register the row's clickable areas

	// "prev" is always the first element in the row, vertically
	// centered, so a point near the left edge at mid-height reliably
	// lands on it regardless of exact font-driven button sizes.
	press(f32.Pt(8, 30))
	frame()
	release(f32.Pt(8, 30))
	frame()

	if p.Current() != 1 {
		t.Fatalf("after clicking prev at page 2: Current() = %d, want 1", p.Current())
	}
	if changedTo != 1 {
		t.Errorf("OnChange fired with %d, want 1", changedTo)
	}
}

func TestPagination_GoToIgnoresOutOfRangeAndNoOp(t *testing.T) {
	calls := 0
	p := New(3).OnChange(func(int) { calls++ })

	p.goTo(0)
	p.goTo(4)
	p.goTo(1) // already current
	if calls != 0 {
		t.Errorf("OnChange fired %d times for out-of-range/no-op moves, want 0", calls)
	}
	if p.Current() != 1 {
		t.Errorf("Current() = %d after no-op moves, want 1", p.Current())
	}
}

func TestPagination_PageButtonsCachedAcrossFrames(t *testing.T) {
	p := New(10)
	b1 := p.pageButton(5)
	b2 := p.pageButton(5)
	if b1 != b2 {
		t.Error("pageButton(5) returned different instances on two calls, want the same cached Button")
	}
}
