// Package uitest drives Giode views in tests the way the app's frame
// loop does: it lays the view out with popups and the semantic tree,
// feeds pointer and key input through a real Gio router, and finds
// nodes by kind and label.
package uitest

import (
	"image"
	"time"

	"gioui.org/f32"
	"gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/fonts"
)

func init() {
	fonts.UseGofont()
}

// Harness lays out View in a window of Size.
type Harness struct {
	View func() elements.Element
	Size image.Point
	Now  time.Time

	router input.Router
	nodes  []elements.SemanticNode
}

// New returns a harness with an 800x600 window.
func New(view func() elements.Element) *Harness {
	return &Harness{View: view, Size: image.Pt(800, 600), Now: time.Unix(0, 0)}
}

// Frame lays out one frame and returns its semantic nodes.
func (h *Harness) Frame() []elements.SemanticNode {
	var ops op.Ops
	h.Now = h.Now.Add(16 * time.Millisecond)
	gtx := layout.Context{Ops: &ops, Now: h.Now, Constraints: layout.Exact(h.Size), Source: h.router.Source()}
	elements.BeginFrame()
	elements.BeginSemantics()
	elements.SetWindowSize(gtx.Constraints.Max)
	h.View().Layout(gtx)
	elements.DrainPortals(gtx)
	elements.RunAfterLayout()
	h.nodes = elements.EndSemantics()
	elements.EndFrame()
	h.router.Frame(&ops)
	return h.nodes
}

// Frames lays out n frames.
func (h *Harness) Frames(n int) []elements.SemanticNode {
	for i := 0; i < n; i++ {
		h.Frame()
	}
	return h.nodes
}

// Move moves the pointer to p and lays out two frames.
func (h *Harness) Move(p image.Point) {
	h.router.Queue(pointer.Event{Kind: pointer.Move, Source: pointer.Mouse, Position: f32.Pt(float32(p.X), float32(p.Y))})
	h.Frames(2)
}

// Click clicks at p and lays out the frames a real click produces.
func (h *Harness) Click(p image.Point) {
	pos := f32.Pt(float32(p.X), float32(p.Y))
	h.router.Queue(pointer.Event{Kind: pointer.Move, Source: pointer.Mouse, Position: pos})
	h.router.Queue(pointer.Event{Kind: pointer.Press, Source: pointer.Mouse, Position: pos, Buttons: pointer.ButtonPrimary})
	h.Frame()
	h.router.Queue(pointer.Event{Kind: pointer.Release, Source: pointer.Mouse, Position: pos})
	h.Frames(2)
}

// Key presses and releases a key, then lays out two frames.
func (h *Harness) Key(name key.Name) {
	h.router.Queue(key.Event{Name: name, State: key.Press}, key.Event{Name: name, State: key.Release})
	h.Frames(2)
}

// Find returns the first node of the last frame with kind and label.
func (h *Harness) Find(kind elements.SemanticKind, label string) (elements.SemanticNode, bool) {
	for _, n := range h.nodes {
		if n.Kind == kind && n.Label == label {
			return n, true
		}
	}
	return elements.SemanticNode{}, false
}
