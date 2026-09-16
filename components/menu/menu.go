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
	submenu   *Menu
}

// Menu is a stateful context menu. Register it with App.SetContextMenu
// to open it with a right-click anywhere in the window, call
// OpenAt/Toggle directly, or wrap a trigger element with Trigger for a
// click-to-open dropdown menu.
//
// When SetContextMenu is used the app registers the right-click area
// and renders the menu on top automatically; otherwise lay Menu out
// anywhere in the view (its panel renders via elements.RegisterPortal,
// so it paints above everything else regardless of tree position).
type Menu struct {
	items  []Item
	open   bool
	pos    f32.Point
	styles styles.Styles

	// openSubmenu is the index of the item whose submenu is showing,
	// or -1. It persists once hover opens a submenu - set again every
	// frame while that item stays the active one - rather than closing
	// the moment the pointer leaves the row, so moving into the
	// submenu's own (separately positioned) panel doesn't flicker it
	// shut; it only changes when a sibling submenu item is hovered, or
	// this menu closes.
	openSubmenu int
	// isSubmenu marks a Menu opened via SubMenu rather than
	// SetContextMenu/OpenAt/Trigger: its visibility is driven entirely
	// by the parent row's hover (openSubmenu above), so it skips its
	// own dismiss-on-outside-click scrim and Escape handling, which
	// would otherwise need to cover the whole window from deep inside
	// the parent panel's own (much smaller) layout constraints.
	isSubmenu bool

	scrim widget.Clickable
}

// New returns an empty menu. The styles argument is optional:
// Background, BorderRadius, Padding and Color (item text) are
// honored; the rest fall back to defaults.
func New(st ...styles.Styles) *Menu {
	m := &Menu{openSubmenu: -1}
	if len(st) > 0 {
		m.styles = st[0]
	}
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

// SubMenu appends an item that reveals submenu as a nested panel,
// cascading to the side, when hovered. submenu is a *Menu of its own
// - built the same way as any other (Item, IconItem, Separator,
// SubMenu for further nesting) - but must not be registered with
// SetContextMenu, OpenAt, Toggle or Trigger itself: this item's hover
// state is its only way of opening. It returns m for chaining.
func (m *Menu) SubMenu(label string, submenu *Menu) *Menu {
	submenu.isSubmenu = true
	m.items = append(m.items, Item{
		label:   label,
		click:   new(widget.Clickable),
		submenu: submenu,
	})
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
		m.Close()
		return
	}
	m.OpenAt(p)
}

// Close closes the menu, and any open submenu (recursively, in case a
// submenu itself has a submenu open).
func (m *Menu) Close() {
	m.open = false
	m.openSubmenu = -1
	for i := range m.items {
		if sub := m.items[i].submenu; sub != nil {
			sub.Close()
		}
	}
}

// Opened reports whether the menu is open.
func (m *Menu) Opened() bool {
	return m.open
}
