// Package tabs provides a stateful tab bar with per-tab views.
package tabs

import (
	"image/color"

	"gioui.org/widget"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
)

var (
	defaultActive     = color.NRGBA{R: 0x3b, G: 0x82, B: 0xf6, A: 0xff}
	defaultText       = color.NRGBA{R: 0x94, G: 0xa3, B: 0xb8, A: 0xff}
	defaultActiveText = color.NRGBA{R: 0xf8, G: 0xfa, B: 0xfc, A: 0xff}
)

// Tab is one tab: its label and the view shown while selected.
type Tab struct {
	Label string
	View  func() elements.Element
}

// Tabs is a stateful tab bar. Create it once with New and lay it out
// every frame; only the selected tab's view is built.
type Tabs struct {
	tabs     []Tab
	selected int
	clicks   []*widget.Clickable
	onChange func(index int)
	st       styles.Styles
}

// New returns a Tabs with the given tabs.
func New(tabs ...Tab) *Tabs {
	t := &Tabs{tabs: tabs}
	for range tabs {
		t.clicks = append(t.clicks, new(widget.Clickable))
	}
	return t
}

// Selected returns the index of the active tab.
func (t *Tabs) Selected() int {
	return t.selected
}

// SetSelected switches the active tab without invoking OnChange.
func (t *Tabs) SetSelected(i int) {
	if i >= 0 && i < len(t.tabs) {
		t.selected = i
	}
}

// OnChange sets the handler invoked when the active tab changes. It
// returns t for chaining.
func (t *Tabs) OnChange(fn func(index int)) *Tabs {
	t.onChange = fn
	return t
}

// Styles replaces the tab styles: Color is the active indicator and
// BorderColor is the active text. It returns t for chaining.
func (t *Tabs) Styles(st styles.Styles) *Tabs {
	t.st = st
	return t
}
