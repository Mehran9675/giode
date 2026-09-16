package menu

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

func TestMenu_Trigger_ClickTogglesOpen(t *testing.T) {
	elements.SetWindowSize(image.Pt(400, 400))
	defer elements.SetWindowSize(image.Point{})

	clicked := false
	m := New().Item("Item", func() { clicked = true })
	trigger := elements.Box(styles.Styles{Width: 60, Height: 30})
	wrapped := m.Trigger(trigger)

	var ops op.Ops
	var router input.Router
	gtx := layout.Context{Ops: &ops, Source: router.Source(), Constraints: layout.Exact(image.Pt(400, 400))}

	frame := func() {
		ops.Reset()
		wrapped.Layout(gtx)
		m.Layout(gtx)
		elements.DrainPortals(gtx)
		router.Frame(&ops)
	}
	press := func(p f32.Point) {
		router.Queue(pointer.Event{Kind: pointer.Press, Source: pointer.Mouse, Position: p, Buttons: pointer.ButtonPrimary})
	}
	release := func(p f32.Point) {
		router.Queue(pointer.Event{Kind: pointer.Release, Source: pointer.Mouse, Position: p, Buttons: 0})
	}

	frame() // register the trigger's clickable area

	press(f32.Pt(10, 10))
	frame()
	release(f32.Pt(10, 10))
	frame()
	if !m.Opened() {
		t.Fatalf("after clicking the trigger: Opened() = false, want true")
	}

	press(f32.Pt(10, 10))
	frame()
	release(f32.Pt(10, 10))
	frame()
	if m.Opened() {
		t.Fatalf("after clicking the trigger again: Opened() = true, want false")
	}
	if clicked {
		t.Errorf("the menu item's own action fired from clicking the trigger, want it untouched")
	}
}

func TestMenu_SubMenu_HoverOpensAndPersists(t *testing.T) {
	elements.SetWindowSize(image.Pt(400, 400))
	defer elements.SetWindowSize(image.Point{})

	sub := New().Item("Sub item", func() {})
	m := New().SubMenu("More", sub)
	m.OpenAt(f32.Pt(0, 0))

	var ops op.Ops
	var router input.Router
	gtx := layout.Context{Ops: &ops, Source: router.Source(), Constraints: layout.Exact(image.Pt(400, 400))}

	frame := func() {
		ops.Reset()
		m.Layout(gtx)
		elements.DrainPortals(gtx)
		router.Frame(&ops)
	}
	moveTo := func(p f32.Point) {
		router.Queue(pointer.Event{Kind: pointer.Move, Source: pointer.Mouse, Position: p})
	}

	frame() // registers the "More" row's hover area

	// The panel is clamped inside the margin near the (0,0) anchor, so
	// its first (and only) row starts a little inside that corner.
	moveTo(f32.Pt(30, 25))
	frame()
	frame()
	if !sub.Opened() {
		t.Fatalf("after hovering the submenu item: submenu Opened() = false, want true")
	}

	// Moving away from the row (but not onto a sibling submenu item -
	// there isn't one here) must NOT close it: it persists until
	// superseded, so moving into the submenu's own panel doesn't
	// flicker it shut.
	moveTo(f32.Pt(390, 390))
	frame()
	frame()
	if !sub.Opened() {
		t.Fatalf("after moving away (no sibling submenu hovered): submenu Opened() = false, want true (should persist)")
	}

	m.Close()
	if sub.Opened() {
		t.Errorf("after closing the parent menu: submenu Opened() = true, want false (Close is recursive)")
	}
}
