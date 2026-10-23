package elements

import (
	"image"
	"testing"

	"gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
)

// backHarness lays out a frame the way the app does: view first, then
// DispatchBack and RunAfterLayout.
type backHarness struct {
	r    input.Router
	view func()
}

func (h *backHarness) frame() {
	var ops op.Ops
	gtx := layout.Context{Ops: &ops, Constraints: layout.Constraints{Max: image.Pt(400, 400)}, Source: h.r.Source()}
	h.view()
	DispatchBack(gtx)
	RunAfterLayout()
	h.r.Frame(&ops)
}

// press queues a key press and reports whether the router took it,
// which is what Android uses to decide whether Back was handled.
func (h *backHarness) press(name key.Name) bool {
	h.r.Queue(key.Event{Name: name, State: key.Press})
	_, handled := h.r.WakeupTime()
	return handled
}

func TestHandleBack_TopmostWins(t *testing.T) {
	var page, dialog int
	dialogOpen := true
	h := &backHarness{view: func() {
		HandleBack(func() { page++ })
		if dialogOpen {
			HandleDismiss(func() { dialog++; dialogOpen = false })
		}
	}}
	h.frame()

	if !h.press(key.NameBack) {
		t.Fatal("Back not handled with handlers registered")
	}
	h.frame()
	if dialog != 1 || page != 0 {
		t.Fatalf("first Back: dialog=%d page=%d, want 1 0", dialog, page)
	}

	h.press(key.NameBack)
	h.frame()
	if dialog != 1 || page != 1 {
		t.Fatalf("second Back: dialog=%d page=%d, want 1 1", dialog, page)
	}
}

func TestHandleBack_EscapeOnlyForDismiss(t *testing.T) {
	var page int
	h := &backHarness{view: func() { HandleBack(func() { page++ }) }}
	h.frame()
	if h.press(key.NameEscape) {
		t.Error("Escape handled by a Back-only handler")
	}
	h.frame()
	if page != 0 {
		t.Errorf("page = %d, want 0", page)
	}

	var closed int
	h.view = func() {
		HandleBack(func() { page++ })
		HandleDismiss(func() { closed++ })
	}
	h.frame()
	h.press(key.NameEscape)
	h.frame()
	if closed != 1 || page != 0 {
		t.Errorf("Escape: closed=%d page=%d, want 1 0", closed, page)
	}
}

func TestHandleBack_UnhandledWithoutHandlers(t *testing.T) {
	h := &backHarness{view: func() {}}
	h.frame()
	if h.press(key.NameBack) {
		t.Error("Back handled with no handlers, so Android couldn't close the app")
	}
}

func TestHandleBack_MeasureDiscards(t *testing.T) {
	var ran int
	h := &backHarness{view: func() {
		var ops op.Ops
		gtx := layout.Context{Ops: &ops, Constraints: layout.Constraints{Max: image.Pt(10, 10)}}
		Measure(gtx, Raw(func(gtx layout.Context) layout.Dimensions {
			HandleDismiss(func() { ran++ })
			return layout.Dimensions{}
		}))
	}}
	h.frame()
	if h.press(key.NameBack) {
		t.Error("Back handled by a handler registered while measuring")
	}
}
