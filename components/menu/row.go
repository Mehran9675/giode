package menu

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/components/icon"
	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// row returns the element for one menu row. index is the item's
// position in m.items, needed to track which row's submenu (if any)
// is currently open.
func (m *Menu) row(item *Item, index, width int, st styles.Styles) elements.Element {
	if item.separator {
		// Divider always fills its ambient Max.X rather than hugging
		// (that's its point when used standalone), so it must be
		// wrapped to the shared row width explicitly once known -
		// otherwise it would stretch to whatever loose constraint the
		// panel happens to have been given.
		divider := elements.Divider(styles.Styles{Color: properties.CalcColorReverse(defaultSeparator), Height: 1})
		if width > 0 {
			return elements.Box(styles.Styles{Width: properties.Width(width)}, divider)
		}
		return divider
	}
	if item.disabled {
		st := styles.Styles{
			Padding: properties.SymmetricInset(28, 8),
			Color:   properties.CalcColorReverse(defaultDisabled),
		}
		if width > 0 {
			st.Width = properties.Width(width)
		}
		return elements.Box(st, elements.Text(item.label))
	}
	return &menuRow{
		menu:  m,
		item:  item,
		index: index,
		width: width,
		st:    st,
	}
}

// menuRow is a hoverable, clickable item row - or, when item.submenu
// is set, a hoverable row that opens a cascading nested panel instead
// of firing an action on click.
type menuRow struct {
	menu  *Menu
	item  *Item
	index int
	width int
	st    styles.Styles
}

// Layout lays out the row, handles its click, and - for a row with a
// submenu - opens/renders it while it's the active one (see
// Menu.openSubmenu).
func (r *menuRow) Layout(gtx layout.Context) layout.Dimensions {
	item := r.item
	if item.submenu == nil && item.click.Clicked(gtx) {
		if item.action != nil {
			item.action()
		}
		r.menu.Close()
		gtx.Execute(op.InvalidateCmd{})
	}
	dims := item.click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		st := styles.Styles{
			Padding:      properties.SymmetricInset(28, 8),
			BorderRadius: 6,
		}
		// Force the row to the shared width via the Width style, not
		// by poking the incoming constraints directly: an unset Width
		// makes this box hug, which zeroes the constraint minimum
		// regardless of what a caller set it to, discarding a
		// raw-constraint override silently. Going through Width keeps
		// the box's own "am I fixed-size?" check (hugMain) in sync
		// with what's actually being asked of it.
		if r.width > 0 {
			st.Width = properties.Width(r.width)
		}
		hovered := item.click.Hovered()
		if hovered {
			st.Background = properties.CalcColorReverse(defaultHover)
			pointer.CursorPointer.Add(gtx.Ops)
			if item.submenu != nil {
				r.menu.openSubmenu = r.index
			}
		}
		var children []elements.Element
		if item.icon != "" {
			children = append(children,
				icon.Material(item.icon, styles.Styles{Color: r.st.Color, Width: 16, Height: 16}),
				elements.Box(styles.Styles{Width: 8}),
			)
		}
		children = append(children, elements.Box(styles.Styles{FlexGrow: 1}, elements.Text(item.label, styles.Styles{Color: r.st.Color})))
		if item.submenu != nil {
			children = append(children,
				elements.Box(styles.Styles{Width: 8}),
				icon.Material("chevron-right", styles.Styles{Color: r.st.Color, Width: 14, Height: 14}),
			)
		}
		return elements.Box(st, elements.Row(styles.Styles{Align: properties.AlignCenter}, children...)).Layout(gtx)
	})
	// Keep the item clickable even when a frameless window's
	// whole-window drag area happens to reach this portal - see
	// button.Layout's identical comment.
	kit.BlockDrag(gtx, dims.Size)

	if item.submenu != nil && r.menu.openSubmenu == r.index {
		anchor := elements.AbsolutePosition().Add(image.Pt(dims.Size.X, 0))
		item.submenu.OpenAt(f32.Pt(float32(anchor.X), float32(anchor.Y)))
		item.submenu.Layout(gtx)
	}
	return dims
}
