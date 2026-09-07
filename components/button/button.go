// Package button provides a stateful clickable button component built
// on the stateless elements.
package button

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/mehran9675/giode/styles"
)

// Button is a stateful clickable button. Create it once with New and
// lay it out every frame; label, styles and handler may change at any
// time.
type Button struct {
	label     string
	st        styles.Styles
	onClick   func()
	clickable widget.Clickable
}

// New returns a Button with the given label.
func New(label string) *Button {
	return &Button{label: label}
}

// OnClick sets the handler invoked on every click. It returns b for
// chaining.
func (b *Button) OnClick(fn func()) *Button {
	b.onClick = fn
	return b
}

// Styles replaces the button styles. It returns b for chaining.
func (b *Button) Styles(st styles.Styles) *Button {
	b.st = st
	return b
}

// Label replaces the button label. It returns b for chaining.
func (b *Button) Label(label string) *Button {
	b.label = label
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
	}
	return b.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return b.layoutVisual(gtx)
	})
}
