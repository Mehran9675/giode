package menu

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/components/icon"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// row returns the element for one menu row.
func (m *Menu) row(item *Item, width int, st styles.Styles) elements.Element {
	if item.separator {
		return elements.Divider(styles.Styles{Color: properties.CalcColorReverse(defaultSeparator), Height: 1})
	}
	if item.disabled {
		return elements.Box(styles.Styles{
			Padding: properties.SymmetricInset(28, 8),
			Color:   properties.CalcColorReverse(defaultDisabled),
		}, elements.Text(item.label))
	}
	return &menuRow{
		menu:  m,
		item:  item,
		width: width,
		st:    st,
	}
}

// menuRow is a hoverable, clickable item row.
type menuRow struct {
	menu  *Menu
	item  *Item
	width int
	st    styles.Styles
}

// Layout lays out the row and handles its click.
func (r *menuRow) Layout(gtx layout.Context) layout.Dimensions {
	item := r.item
	if item.click.Clicked(gtx) {
		if item.action != nil {
			item.action()
		}
		r.menu.open = false
		gtx.Execute(op.InvalidateCmd{})
	}
	return item.click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if r.width > 0 {
			gtx.Constraints.Min.X = r.width
			gtx.Constraints.Max.X = r.width
		}
		st := styles.Styles{
			Padding:      properties.SymmetricInset(28, 8),
			BorderRadius: 6,
		}
		if item.click.Hovered() {
			st.Background = properties.CalcColorReverse(defaultHover)
			pointer.CursorPointer.Add(gtx.Ops)
		}
		var children []elements.Element
		if item.icon != "" {
			children = append(children,
				icon.Material(item.icon, styles.Styles{Color: r.st.Color, Width: 16, Height: 16}),
				elements.Box(styles.Styles{Width: 8}),
			)
		}
		children = append(children, elements.Text(item.label, styles.Styles{Color: r.st.Color}))
		return elements.Box(st, elements.Row(styles.Styles{Align: properties.AlignCenter}, children...)).Layout(gtx)
	})
}
