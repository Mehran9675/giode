package elements

import (
	"image"
	"testing"
	"time"

	"gioui.org/f32"
	"gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/styles"
)

// harness lays out an element against a real Gio router, so tests can
// queue pointer and key events.
type harness struct {
	r   input.Router
	el  Element
	now time.Time
}

func (h *harness) frame() {
	var ops op.Ops
	h.now = h.now.Add(16 * time.Millisecond)
	gtx := layout.Context{Ops: &ops, Now: h.now, Constraints: layout.Constraints{Max: image.Pt(400, 400)}, Source: h.r.Source()}
	h.el.Layout(gtx)
	h.r.Frame(&ops)
}

func (h *harness) click(x, y float32, btn pointer.Buttons, mods key.Modifiers) {
	pos := f32.Pt(x, y)
	h.r.Queue(pointer.Event{Kind: pointer.Move, Source: pointer.Mouse, Position: pos})
	h.r.Queue(pointer.Event{Kind: pointer.Press, Source: pointer.Mouse, Position: pos, Buttons: btn, Modifiers: mods})
	h.frame()
	h.r.Queue(pointer.Event{Kind: pointer.Release, Source: pointer.Mouse, Position: pos, Modifiers: mods})
	h.frame()
}

func sized(w, h int) Element {
	return Raw(func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: image.Pt(w, h)}
	})
}

func TestPressable_ClickKinds(t *testing.T) {
	var events []ClickEvent
	var doubles, rights int
	p := NewPressable(func(PressState) Element { return sized(100, 40) }).
		OnClick(func(e ClickEvent) { events = append(events, e) }).
		OnDoubleClick(func(ClickEvent) { doubles++ }).
		OnSecondaryClick(func(ClickEvent) { rights++ })
	h := &harness{el: p, now: time.Unix(0, 0)}
	h.frame()

	h.click(10, 10, pointer.ButtonPrimary, key.ModShift)
	h.click(10, 10, pointer.ButtonPrimary, 0)
	if len(events) != 2 || events[0].Clicks != 1 || events[1].Clicks != 2 || doubles != 1 {
		t.Fatalf("events %+v, doubles %d: want a single then a double click", events, doubles)
	}
	if !events[0].Shift() || events[0].Source != ClickPointer || events[0].Position != image.Pt(10, 10) {
		t.Errorf("first event = %+v", events[0])
	}
	h.now = h.now.Add(time.Second)
	h.click(10, 10, pointer.ButtonSecondary, 0)
	if rights != 1 || events[2].Button != MouseSecondary || events[2].Clicks != 1 {
		t.Errorf("right click: %d calls, event %+v", rights, events[2])
	}
	h.click(300, 300, pointer.ButtonPrimary, 0)
	if len(events) != 3 {
		t.Error("click outside the content fired")
	}
}

func TestPressable_ChildrenKeepTheirClicks(t *testing.T) {
	var outer, inner int
	child := NewPressable(func(PressState) Element { return sized(20, 20) }).OnClick(func(ClickEvent) { inner++ })
	p := NewPressable(func(PressState) Element {
		return Row(styles.Styles{}, child, sized(80, 20))
	}).OnClick(func(ClickEvent) { outer++ })
	h := &harness{el: p, now: time.Unix(0, 0)}
	h.frame()
	h.click(5, 5, pointer.ButtonPrimary, 0)
	if inner != 1 || outer != 0 {
		t.Errorf("click on child: inner %d, outer %d - want the child only", inner, outer)
	}
	h.now = h.now.Add(time.Second)
	h.click(50, 5, pointer.ButtonPrimary, 0)
	if outer != 1 {
		t.Errorf("click beside the child: outer %d, want 1", outer)
	}
}

func TestPressable_KeyboardAndState(t *testing.T) {
	var events []ClickEvent
	var states []PressState
	p := NewPressable(func(s PressState) Element {
		states = append(states, s)
		return sized(100, 40)
	}).OnClick(func(e ClickEvent) { events = append(events, e) })
	h := &harness{el: p, now: time.Unix(0, 0)}
	h.frame()
	h.r.Queue(pointer.Event{Kind: pointer.Move, Source: pointer.Mouse, Position: f32.Pt(5, 5)})
	h.frame()
	h.frame()
	if !states[len(states)-1].Hovered {
		t.Error("content not told it's hovered")
	}
	h.r.Source().Execute(key.FocusCmd{Tag: p})
	h.frame()
	h.r.Queue(key.Event{Name: key.NameSpace, State: key.Press}, key.Event{Name: key.NameSpace, State: key.Release})
	h.frame()
	if len(events) != 1 || events[0].Source != ClickKeyboard {
		t.Fatalf("Space on focused pressable: %+v", events)
	}
	if !p.Focused() {
		t.Error("Focused() false after focusing")
	}
}

func TestPressable_Disabled(t *testing.T) {
	clicks := 0
	var last PressState
	p := NewPressable(func(s PressState) Element {
		last = s
		return sized(100, 40)
	}).OnClick(func(ClickEvent) { clicks++ }).Disabled(true)
	h := &harness{el: p, now: time.Unix(0, 0)}
	h.frame()
	h.click(10, 10, pointer.ButtonPrimary, 0)
	p.Click()
	if clicks != 0 || !last.Disabled {
		t.Errorf("disabled: %d clicks, state %+v", clicks, last)
	}
	p.Disabled(false).Click()
	if clicks != 1 {
		t.Error("programmatic click after enabling didn't fire")
	}
}

func TestBox_OnClick(t *testing.T) {
	var got ClickEvent
	b := Box(styles.Styles{Width: 100, Height: 40}).OnClick(func(e ClickEvent) { got = e })
	h := &harness{el: b, now: time.Unix(0, 0)}
	h.frame()
	h.click(10, 10, pointer.ButtonPrimary, 0)
	if got.Clicks != 1 {
		t.Errorf("Box.OnClick didn't fire: %+v", got)
	}
}
