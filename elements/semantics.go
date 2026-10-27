package elements

import (
	"fmt"
	"image"
	"sync"

	"gioui.org/layout"
)

// SemanticKind classifies a node of the semantic tree.
type SemanticKind string

const (
	SemanticText     SemanticKind = "text"
	SemanticButton   SemanticKind = "button"
	SemanticInput    SemanticKind = "input"
	SemanticCheckbox SemanticKind = "checkbox"
	SemanticSlider   SemanticKind = "slider"
	SemanticSelect   SemanticKind = "select"
	SemanticTab      SemanticKind = "tab"
	SemanticMenuItem SemanticKind = "menuitem"
	SemanticLink     SemanticKind = "link"
	SemanticImage    SemanticKind = "image"
	SemanticScroll   SemanticKind = "scroll"
	SemanticCustom   SemanticKind = "custom"
)

// SemanticNode describes one meaningful piece of the UI - a text run,
// a button, a field - as laid out in the last frame. The semantic tree
// is what automation (the built-in MCP server, tests) sees of the UI.
type SemanticNode struct {
	// ID identifies the node. Nodes with a Target keep their ID for
	// the target's lifetime; others are numbered in layout order.
	ID string `json:"id"`
	// Kind classifies the node.
	Kind SemanticKind `json:"kind"`
	// Label is the visible text or accessible name.
	Label string `json:"label,omitempty"`
	// Value is the current value of a control (field text, checked
	// state, slider position).
	Value string `json:"value,omitempty"`
	// Bounds is the node's window-relative area, in pixels, including
	// for children of rows, stacks and scrolled content.
	Bounds image.Rectangle `json:"bounds"`
	// Disabled marks a control that doesn't accept input.
	Disabled bool `json:"disabled,omitempty"`
	// Scroll is the ID of the innermost scroll container the node is
	// laid out in, or "" outside any. A node in a scroll container may
	// lie outside the container's Bounds, scrolled out of view.
	Scroll string `json:"scroll,omitempty"`
	// Target is the stateful component behind the node, if any, so
	// automation can act on it (e.g. call Click). It is not
	// serialized.
	Target any `json:"-"`
}

// semantics is the recorder for the frame being laid out; nil (the
// default) disables recording, which then costs a nil check.
var (
	semantics   *[]SemanticNode
	scrollStack []string
	semanticIDs = map[any]string{}
	semanticSeq = map[SemanticKind]int{}
	semanticMu  sync.Mutex
)

// BeginSemantics starts recording the semantic tree for a frame. The
// app calls it when semantics are enabled (the MCP server enables
// them); component authors don't need to.
func BeginSemantics() {
	nodes := make([]SemanticNode, 0, 64)
	semantics = &nodes
	scrollStack = scrollStack[:0]
}

// EndSemantics stops recording and returns the frame's nodes.
func EndSemantics() []SemanticNode {
	if semantics == nil {
		return nil
	}
	nodes := *semantics
	semantics = nil
	return nodes
}

// SemanticsEnabled reports whether the current frame records a
// semantic tree. Components can skip computing labels when it doesn't.
func SemanticsEnabled() bool {
	return semantics != nil
}

// RecordSemantic adds a node to the semantic tree of the frame being
// laid out, a no-op unless recording. Call it from a custom
// component's Layout, after laying out, with the component's size:
//
//	elements.RecordSemantic(gtx, elements.SemanticNode{
//		Kind: elements.SemanticButton, Label: b.label, Target: b,
//	}, dims.Size)
//
// Bounds and ID are filled in when left empty.
func RecordSemantic(gtx layout.Context, n SemanticNode, size image.Point) {
	if semantics == nil {
		return
	}
	if n.Bounds.Empty() {
		origin := AbsolutePosition()
		n.Bounds = image.Rectangle{Min: origin, Max: origin.Add(size)}
	}
	if n.ID == "" {
		n.ID = semanticID(n.Kind, n.Target)
	}
	if n.Scroll == "" && len(scrollStack) > 0 {
		n.Scroll = scrollStack[len(scrollStack)-1]
	}
	*semantics = append(*semantics, n)
}

// PushSemanticScroll marks the nodes recorded until pop is called as
// content of the scroll container target, which should then record
// itself as a SemanticScroll node with that Target. A scroll component
// calls it around laying out its content, so automation can scroll the
// container to reach a node:
//
//	pop := elements.PushSemanticScroll(s)
//	dims := s.content().Layout(cgtx)
//	pop()
//	...
//	elements.RecordSemantic(gtx, elements.SemanticNode{Kind: elements.SemanticScroll, Target: s}, viewport)
//
// For automation to scroll it, target should have the methods
// ScrollBy(dx, dy float32) and ScrollOffset() (x, y float32), in
// pixels.
func PushSemanticScroll(target any) (pop func()) {
	if semantics == nil {
		return func() {}
	}
	scrollStack = append(scrollStack, semanticID(SemanticScroll, target))
	n := len(scrollStack)
	return func() {
		if len(scrollStack) >= n {
			scrollStack = scrollStack[:n-1]
		}
	}
}

// recordSemantic records a node without a target.
func recordSemantic(gtx layout.Context, kind SemanticKind, label, value string, size image.Point) {
	if semantics == nil || label == "" {
		return
	}
	RecordSemantic(gtx, SemanticNode{Kind: kind, Label: label, Value: value}, size)
}

// semanticID returns a stable ID for target, or a frame-order ID for
// nodes without one.
func semanticID(kind SemanticKind, target any) string {
	semanticMu.Lock()
	defer semanticMu.Unlock()
	if target != nil {
		if id, ok := semanticIDs[target]; ok {
			return id
		}
	}
	semanticSeq[kind]++
	id := fmt.Sprintf("%s-%d", kind, semanticSeq[kind])
	if target != nil {
		semanticIDs[target] = id
	} else {
		// Untargeted nodes are renumbered every frame.
		id = fmt.Sprintf("%s@%d", kind, len(*semantics))
		semanticSeq[kind]--
	}
	return id
}
