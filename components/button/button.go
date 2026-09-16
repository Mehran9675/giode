// Package button provides a stateful clickable button component built
// on the stateless elements.
package button

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"

	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/styles"
)

// Variant selects the button's visual style.
type Variant string

const (
	// VariantContained fills the button with Background (or black if
	// unset) at all times - the default, and today's only look before
	// Variant existed.
	VariantContained Variant = "contained"
	// VariantOutlined is transparent with a visible border, filling
	// faintly with the border color on hover/press.
	VariantOutlined Variant = "outlined"
	// VariantGhost is fully transparent at rest - only the label is
	// visible - filling faintly with Background on hover/press.
	VariantGhost Variant = "ghost"
)

// Button is a stateful clickable button. Create it once with New and
// lay it out every frame; label, styles and handler may change at any
// time.
type Button struct {
	label     string
	st        styles.Styles
	variant   Variant
	onClick   func()
	clickable widget.Clickable
}

// New returns a Button with the given label. The styles argument is
// optional.
func New(label string, st ...styles.Styles) *Button {
	b := &Button{label: label}
	if len(st) > 0 {
		b.st = st[0]
	}
	return b
}

// OnClick sets the handler invoked on every click. It returns b for
// chaining.
func (b *Button) OnClick(fn func()) *Button {
	b.onClick = fn
	return b
}

// Label replaces the button label. It returns b for chaining.
func (b *Button) Label(label string) *Button {
	b.label = label
	return b
}

// Variant sets the button's visual style (VariantContained, the
// default, VariantOutlined or VariantGhost). It returns b for
// chaining.
func (b *Button) Variant(v Variant) *Button {
	b.variant = v
	return b
}

// Click simulates a click programmatically.
func (b *Button) Click() {
	b.clickable.Click()
}

// Layout lays out and updates the button.
func (b *Button) Layout(gtx layout.Context) layout.Dimensions {
	if b.clickable.Clicked(gtx) && b.onClick != nil {
		b.onClick()
		// The handler may have changed state that this frame's view
		// already read before this click was processed (e.g. a
		// counter formatted into a Text earlier in the tree); request
		// another frame so the change shows up right away instead of
		// waiting for an unrelated event to trigger a repaint.
		gtx.Execute(op.InvalidateCmd{})
	}
	dims := b.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return b.layoutVisual(gtx)
	})
	// Keep the button clickable even when an ancestor is a frameless
	// window's whole-window drag area (see the root package's
	// layoutDrag): without this, pressing the button would be
	// hijacked into a native window-move instead of a click.
	kit.BlockDrag(gtx, dims.Size)
	return dims
}
