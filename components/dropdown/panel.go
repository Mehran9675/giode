package dropdown

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// dropdown lays out the option list and processes item clicks.
func (s *Select) dropdown(gtx layout.Context, st styles.Styles, width int) layout.Dimensions {
	gtx.Constraints.Min.X = width
	gtx.Constraints.Max.X = width
	st.Display = properties.Flex
	st.Direction = properties.Column
	st.Padding = properties.UniformInset(4)

	children := make([]elements.Element, 0, len(s.options))
	for i := range s.options {
		children = append(children, s.optionRow(i, st))
	}
	return elements.Box(st, children...).Layout(gtx)
}

// optionRow renders one hoverable option.
func (s *Select) optionRow(i int, st styles.Styles) elements.Element {
	click := s.itemClick[i]
	return elements.Raw(func(gtx layout.Context) layout.Dimensions {
		if click.Clicked(gtx) {
			s.selected = i
			s.open = false
			if s.onChange != nil {
				s.onChange(i)
			}
		}
		return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			rowSt := styles.Styles{
				Padding:      properties.SymmetricInset(12, 8),
				BorderRadius: 4,
			}
			if click.Hovered() {
				rowSt.Background = defaultHover
				pointer.CursorPointer.Add(gtx.Ops)
			}
			return elements.Box(rowSt, elements.Text(s.options[i], styles.Styles{Color: st.Color})).Layout(gtx)
		})
	})
}
