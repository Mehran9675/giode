package dropdown

import (
	"image"
	"testing"

	"gioui.org/f32"
	"gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/fonts"
	"github.com/mehran9675/giode/styles"
)

func init() {
	fonts.UseGofont()
}

type harness struct {
	r    input.Router
	view func() elements.Element
}

// frame lays the view out the way the app does, popups included, and
// returns the semantic nodes.
func (h *harness) frame() []elements.SemanticNode {
	var ops op.Ops
	gtx := layout.Context{Ops: &ops, Constraints: layout.Exact(image.Pt(800, 600)), Source: h.r.Source()}
	elements.BeginFrame()
	elements.BeginSemantics()
	elements.SetWindowSize(gtx.Constraints.Max)
	h.view().Layout(gtx)
	elements.DrainPortals(gtx)
	nodes := elements.EndSemantics()
	elements.EndFrame()
	h.r.Frame(&ops)
	return nodes
}

func (h *harness) click(p image.Point) {
	pos := f32.Pt(float32(p.X), float32(p.Y))
	h.r.Queue(pointer.Event{Kind: pointer.Move, Source: pointer.Mouse, Position: pos})
	h.r.Queue(pointer.Event{Kind: pointer.Press, Source: pointer.Mouse, Position: pos, Buttons: pointer.ButtonPrimary})
	h.frame()
	h.r.Queue(pointer.Event{Kind: pointer.Release, Source: pointer.Mouse, Position: pos})
	h.frame()
	h.frame()
}

func find(nodes []elements.SemanticNode, kind elements.SemanticKind, label string) (elements.SemanticNode, bool) {
	for _, n := range nodes {
		if n.Kind == kind && n.Label == label {
			return n, true
		}
	}
	return elements.SemanticNode{}, false
}

// nested places the select as the second child of a Row inside a Stack,
// below a header - the layout where the list used to open at the
// window's top-left corner.
func nested(s *Select) func() elements.Element {
	return func() elements.Element {
		return elements.Stack(styles.Styles{Gap: 20},
			elements.Box(styles.Styles{Height: 100}),
			elements.Row(styles.Styles{Gap: 16, Wrap: "nowrap"},
				elements.Box(styles.Styles{Width: 150, Height: 20}),
				s,
			),
		)
	}
}

func TestListOpensBelowFieldInsideRowAndStack(t *testing.T) {
	s := New(styles.Styles{Width: 120}, "Small", "Medium", "Large")
	h := &harness{view: nested(s)}
	nodes := h.frame()
	field, ok := find(nodes, elements.SemanticSelect, "Small")
	if !ok {
		t.Fatal("select field not in the semantic tree")
	}
	if field.Bounds.Min.X != 166 || field.Bounds.Min.Y != 120 {
		t.Fatalf("field at %v, want it after the 150dp box and 16dp gap, below the 100dp header and gap", field.Bounds)
	}

	h.click(field.Bounds.Min.Add(image.Pt(10, 10)))
	nodes = h.frame()
	opt, ok := find(nodes, elements.SemanticText, "Medium")
	if !ok {
		t.Fatal("list didn't open")
	}
	if opt.Bounds.Min.X < field.Bounds.Min.X || opt.Bounds.Min.X > field.Bounds.Min.X+20 {
		t.Errorf("option x = %d, want just inside the field's left edge %d", opt.Bounds.Min.X, field.Bounds.Min.X)
	}
	if opt.Bounds.Min.Y < field.Bounds.Max.Y || opt.Bounds.Min.Y > field.Bounds.Max.Y+80 {
		t.Errorf("option y = %d, want just below the field's bottom %d", opt.Bounds.Min.Y, field.Bounds.Max.Y)
	}
}

func TestOutsideClickCloses(t *testing.T) {
	s := New(styles.Styles{Width: 120}, "Small", "Medium")
	h := &harness{view: nested(s)}
	field, _ := find(h.frame(), elements.SemanticSelect, "Small")
	h.click(field.Bounds.Min.Add(image.Pt(10, 10)))
	if !s.Opened() {
		t.Fatal("list didn't open")
	}
	h.click(image.Pt(700, 500))
	if s.Opened() {
		t.Error("click outside the list didn't close it")
	}
	if s.Selected() != 0 {
		t.Error("outside click changed the selection")
	}
}

func TestKeyboard(t *testing.T) {
	changed := -1
	s := New(styles.Styles{Width: 120}, "Small", "Medium", "Large").OnChange(func(i int) { changed = i })
	h := &harness{view: nested(s)}
	field, _ := find(h.frame(), elements.SemanticSelect, "Small")
	h.click(field.Bounds.Min.Add(image.Pt(10, 10)))
	press := func(name key.Name) {
		h.r.Queue(key.Event{Name: name, State: key.Press}, key.Event{Name: name, State: key.Release})
		h.frame()
		h.frame()
	}
	press(key.NameDownArrow)
	press(key.NameDownArrow)
	press(key.NameReturn)
	if s.Opened() || s.Selected() != 2 || changed != 2 {
		t.Fatalf("Down Down Enter: open %v, selected %d, OnChange %d - want closed on Large", s.Opened(), s.Selected(), changed)
	}
	// Enter on the focused, closed field opens it; Escape closes it
	// without changing anything.
	press(key.NameReturn)
	if !s.Opened() {
		t.Fatal("Enter on the focused field didn't open the list")
	}
	press(key.NameUpArrow)
	press(key.NameEscape)
	if s.Opened() || s.Selected() != 2 {
		t.Errorf("Up Escape: open %v, selected %d - want closed, still Large", s.Opened(), s.Selected())
	}
}
