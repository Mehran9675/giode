// Package menu provides a right-click context menu that opens at the
// cursor and keeps itself within the bounds of the window, flipping
// above the cursor when it would overflow the bottom edge.
package menu

import (
	"gioui.org/f32"
	"gioui.org/widget"

	"github.com/mehran9675/giode/styles"
)

// Item is a single row of a menu.
type Item struct {
	label     string
	icon      string
	disabled  bool
	separator bool
	click     *widget.Clickable
	action    func()
}

// Menu is a stateful context menu. Register it with App.SetContextMenu
// to open it with a right-click anywhere in the window, or call
// OpenAt/Toggle directly.
//
// When SetContextMenu is used the app registers the right-click area
// and renders the menu on top automatically; otherwise lay the menu
// out last in the view.
type Menu struct {
	items  []Item
	open   bool
	pos    f32.Point
	styles styles.Styles

	scrim widget.Clickable
}

// New returns an empty menu.
func New() *Menu {
	return &Menu{}
}

// Styles replaces the panel styles. Background, BorderRadius, Padding
// and Color (item text) are honored; the rest fall back to defaults.
func (m *Menu) Styles(st styles.Styles) *Menu {
	m.styles = st
	return m
}

// Item appends a clickable item. It returns m for chaining.
func (m *Menu) Item(label string, action func()) *Menu {
	m.items = append(m.items, Item{
		label:  label,
		click:  new(widget.Clickable),
		action: action,
	})
	return m
}

// IconItem appends an item with a material icon. It returns m for
// chaining.
func (m *Menu) IconItem(iconName, label string, action func()) *Menu {
	m.items = append(m.items, Item{
		label:  label,
		icon:   iconName,
		click:  new(widget.Clickable),
		action: action,
	})
	return m
}

// DisabledItem appends a grayed-out, non-clickable item. It returns m
// for chaining.
func (m *Menu) DisabledItem(label string) *Menu {
	m.items = append(m.items, Item{label: label, disabled: true})
	return m
}

// Separator appends a horizontal line. It returns m for chaining.
func (m *Menu) Separator() *Menu {
	m.items = append(m.items, Item{separator: true})
	return m
}

// OpenAt opens the menu anchored at p, in window pixel coordinates.
func (m *Menu) OpenAt(p f32.Point) {
	m.pos = p
	m.open = true
}

// Toggle flips the menu state, anchored at p.
func (m *Menu) Toggle(p f32.Point) {
	if m.open {
		m.open = false
		return
	}
	m.OpenAt(p)
}

// Close closes the menu.
func (m *Menu) Close() {
	m.open = false
}

// Opened reports whether the menu is open.
func (m *Menu) Opened() bool {
	return m.open
}
