// Package drawer provides a stateful slide-out side panel.
package drawer

import (
	"image/color"
	"time"

	"gioui.org/widget"

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

var defaultBackground = color.NRGBA{R: 0x1e, G: 0x29, B: 0x3b, A: 0xff}

// Drawer is a stateful side panel with a scrim. Create it once with
// New and lay it out last in the view so it renders above the rest of
// the UI.
type Drawer struct {
	side   Side
	width  int
	st     styles.Styles
	open   bool
	openT  time.Time
	closeT time.Time

	scrim widget.Clickable
}

// New returns a Drawer with the given side, width and panel styles.
func New(side Side, width int, st styles.Styles) *Drawer {
	return &Drawer{side: side, width: width, st: st}
}

// Open opens the drawer.
func (d *Drawer) Open() {
	if d.open {
		return
	}
	d.open = true
	d.openT = time.Now()
}

// Close closes the drawer.
func (d *Drawer) Close() {
	if !d.open {
		return
	}
	d.open = false
	d.closeT = time.Now()
}

// Toggle flips the drawer state.
func (d *Drawer) Toggle() {
	if d.open {
		d.Close()
	} else {
		d.Open()
	}
}

// Opened reports whether the drawer is open.
func (d *Drawer) Opened() bool {
	return d.open
}
