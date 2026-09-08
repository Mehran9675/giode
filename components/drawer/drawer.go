// Package drawer provides a stateful slide-out side panel.
package drawer

import (
	"time"

	"gioui.org/widget"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
)

// Side selects which edge the drawer slides from.
type Side int

const (
	// Left slides in from the left edge.
	Left Side = iota
	// Right slides in from the right edge.
	Right
)

var defaultBackground = "#1e293b"

// Drawer is a side panel with a scrim, sliding in over 200ms. It
// implements elements.Element, so it can be placed directly in the
// tree like any other element — lay it out last (e.g. in Stack) so it
// paints above the rest of the UI. It is a no-op (zero size, nothing
// painted) once fully closed.
//
// Create it once with New, outside the view function.
type Drawer struct {
	open    *bool
	wasOpen bool
	side    Side
	width   int
	content func() elements.Element
	st      styles.Styles
	openT   time.Time
	closeT  time.Time

	scrim widget.Clickable
}

// New returns a Drawer bound to open: it slides in from side whenever
// *open is true, and slides back out when *open becomes false
// (including from clicking the scrim). content builds the panel body
// fresh every frame the drawer is shown, so it can reflect state that
// changed since the drawer opened. The styles argument is optional.
func New(open *bool, side Side, width int, content func() elements.Element, st ...styles.Styles) *Drawer {
	d := &Drawer{open: open, side: side, width: width, content: content}
	if len(st) > 0 {
		d.st = st[0]
	}
	return d
}

// Opened reports whether the drawer is currently open (it may still
// be mid slide-out animation after *open turns false).
func (d *Drawer) Opened() bool {
	return d.open != nil && *d.open
}
