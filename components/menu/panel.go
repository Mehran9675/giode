package menu

import (
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// panel lays out the menu body and processes item clicks.
func (m *Menu) panel(gtx layout.Context) layout.Dimensions {
	st := styles.Merge(styles.Styles{
		FlexDirection: properties.FlexDirectionColumn,
		Background:    properties.CalcColorReverse(defaultBackground),
		BorderRadius:  8,
		Padding:       properties.UniformInset(6),
		Color:         properties.CalcColorReverse(defaultText),
	}, m.styles)
	st.FlexDirection = properties.FlexDirectionColumn

	// Bound the panel to the window.
	win := gtx.Constraints.Max
	gtx.Constraints.Max.X = win.X - 2*kit.Margin
	gtx.Constraints.Max.Y = win.Y - 2*kit.Margin

	width := m.measureWidth(gtx, st)

	children := make([]elements.Element, 0, len(m.items))
	for i := range m.items {
		children = append(children, m.row(&m.items[i], width, st))
	}
	return elements.Box(st, children...).Layout(gtx)
}

// measureWidth lays each row out on a scratch context to find the
// widest one, so all rows share the same width.
func (m *Menu) measureWidth(gtx layout.Context, st styles.Styles) int {
	scratch := gtx
	scratch.Ops = new(op.Ops)
	scratch.Source = gtx.Source.Disabled()
	width := 0
	for i := range m.items {
		dims := m.row(&m.items[i], 0, st).Layout(scratch)
		if dims.Size.X > width {
			width = dims.Size.X
		}
	}
	return width
}
