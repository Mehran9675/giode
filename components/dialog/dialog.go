// Package dialog provides a stateful centered modal dialog.
package dialog

import (
	"image/color"

	"gioui.org/widget"

	"github.com/mehran9675/giode/styles"
)

var defaultBackground = color.NRGBA{R: 0x1e, G: 0x29, B: 0x3b, A: 0xff}

// Dialog is a stateful modal dialog with a dimmed scrim. Create it
// once with New and lay it out last in the view so it renders above
// the rest of the UI.
type Dialog struct {
	open  bool
	st    styles.Styles
	scrim widget.Clickable
}

// New returns a Dialog with the given panel styles.
func New(st styles.Styles) *Dialog {
	return &Dialog{st: st}
}

// Open opens the dialog.
func (d *Dialog) Open() {
	d.open = true
}

// Close closes the dialog.
func (d *Dialog) Close() {
	d.open = false
}

// Opened reports whether the dialog is open.
func (d *Dialog) Opened() bool {
	return d.open
}
