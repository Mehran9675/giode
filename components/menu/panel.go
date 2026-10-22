package menu

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// defaultPanelPadding is the panel's padding in dp when New is given
// none.
const defaultPanelPadding = 6

// panel lays out the menu body and processes item clicks.
func (m *Menu) panel(gtx layout.Context) layout.Dimensions {
	margin := gtx.Dp(unit.Dp(kit.Margin))
	st := styles.Merge(styles.Styles{
		FlexDirection: properties.FlexDirectionColumn,
		Background:    properties.CalcColorReverse(defaultBackground),
		BorderRadius:  8,
		Padding:       properties.UniformInset(defaultPanelPadding),
		Color:         properties.CalcColorReverse(defaultText),
	}, m.styles)
	st.FlexDirection = properties.FlexDirectionColumn

	// Bound the panel to the window and let it hug its content: the
	// incoming gtx is bounded by the window regardless of how deep in
	// the tree Layout was called from - see elements.RegisterPortal's
	// use of elements.SetWindowSize - and is typically exact (Min ==
	// Max, the window size) rather than loose, so the minimum must be
	// cleared or the panel would be forced to fill the window.
	win := gtx.Constraints.Max
	gtx.Constraints.Min = image.Point{}
	gtx.Constraints.Max.X = win.X - 2*margin
	gtx.Constraints.Max.Y = win.Y - 2*margin

	width := m.measureWidth(gtx, st)

	rows := make([]elements.Element, 0, len(m.items))
	for i := range m.items {
		rows = append(rows, m.row(&m.items[i], i, width, st))
	}
	return elements.Box(st, elements.Raw(func(gtx layout.Context) layout.Dimensions {
		return stackRows(gtx, rows)
	})).Layout(gtx)
}

// panelPaddingTop is the panel's top padding in dp: the Padding given
// to New, or the default.
func (m *Menu) panelPaddingTop() int {
	if m.styles.Padding != (properties.Inset{}) {
		return m.styles.Padding.Top
	}
	return defaultPanelPadding
}

// stackRows lays the rows out one under the other. It does the stacking
// itself rather than through a column Box because a flex container
// does not report each child's own position to
// elements.AbsolutePosition, only the container's origin: every row,
// and so every submenu anchored off a row, would see the top of the
// panel as its position.
func stackRows(gtx layout.Context, rows []elements.Element) layout.Dimensions {
	gtx.Constraints.Min.Y = 0
	var size image.Point
	for _, row := range rows {
		cgtx := gtx
		cgtx.Constraints.Max.Y = max(0, gtx.Constraints.Max.Y-size.Y)
		off := image.Pt(0, size.Y)
		trans := op.Offset(off).Push(gtx.Ops)
		pop := elements.PushOffset(off)
		dims := row.Layout(cgtx)
		pop()
		trans.Pop()
		size.X = max(size.X, dims.Size.X)
		size.Y += dims.Size.Y
	}
	return layout.Dimensions{Size: size}
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
		dims := m.row(&m.items[i], i, 0, st).Layout(scratch)
		if dims.Size.X > width {
			width = dims.Size.X
		}
	}
	return width
}
