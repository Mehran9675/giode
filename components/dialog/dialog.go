// Package dialog provides a stateful centered modal dialog.
package dialog

import (
	"gioui.org/widget"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
)

var defaultBackground = "#1e293b"

// Dialog is a modal dialog with a dimmed scrim. It implements
// elements.Element, so it can be placed directly in the tree like any
// other element — lay it out last (e.g. in Stack) so it paints above
// the rest of the UI. It is a no-op (zero size, nothing painted)
// whenever *open is false.
//
// Create it once with New, outside the view function.
type Dialog struct {
	open    *bool
	content func() elements.Element
	st      styles.Styles
	scrim   widget.Clickable
}

// New returns a Dialog bound to open: it is visible whenever *open is
// true, and clicking the scrim sets *open back to false. content
// builds the panel body fresh every frame the dialog is shown, so it
// can reflect state that changed since the dialog opened. The styles
// argument is optional and applies to the panel.
func New(open *bool, content func() elements.Element, st ...styles.Styles) *Dialog {
	d := &Dialog{open: open, content: content}
	if len(st) > 0 {
		d.st = st[0]
	}
	return d
}

// Opened reports whether the dialog is currently visible.
func (d *Dialog) Opened() bool {
	return d.open != nil && *d.open
}
