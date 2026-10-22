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
	// Bounds is the node's window-relative area, in pixels. Like
	// AbsolutePosition it is exact for content placed directly in a
	// Box and approximate for children of multi-child flex containers.
	Bounds image.Rectangle `json:"bounds"`
	// Disabled marks a control that doesn't accept input.
	Disabled bool `json:"disabled,omitempty"`
	// Target is the stateful component behind the node, if any, so
	// automation can act on it (e.g. call Click). It is not
	// serialized.
	Target any `json:"-"`
}

// semantics is the recorder for the frame being laid out; nil (the
// default) disables recording, which then costs a nil check.
var (
	semantics   *[]SemanticNode
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
	*semantics = append(*semantics, n)
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
