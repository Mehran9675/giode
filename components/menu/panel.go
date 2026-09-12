package menu

import (
	"image"

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

	// Bound the panel to the window and let it hug its content: the
	// incoming gtx is the raw root constraint (Menu.Layout is laid out
	// directly off the window, not through a normal flex parent that
	// would otherwise neutralize this), which is typically exact
	// (Min == Max, the window size) rather than loose, so the minimum
	// must be cleared or the panel would be forced to fill the window.
	win := gtx.Constraints.Max
	gtx.Constraints.Min = image.Point{}
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
// widest one, so all rows share the same width. Separators are
// skipped: Divider always fills its ambient constraint rather than
// hugging its (nonexistent) content, so measuring it would report the
// scratch context's width instead of anything meaningful.
func (m *Menu) measureWidth(gtx layout.Context, st styles.Styles) int {
	scratch := gtx
	scratch.Ops = new(op.Ops)
	scratch.Source = gtx.Source.Disabled()
	width := 0
	for i := range m.items {
		if m.items[i].separator {
			continue
		}
		dims := m.row(&m.items[i], 0, st).Layout(scratch)
		if dims.Size.X > width {
			width = dims.Size.X
		}
	}
	return width
}
