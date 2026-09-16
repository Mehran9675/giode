package dropdown

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"

	"github.com/mehran9675/giode/components/icon"
	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// Layout lays out the field and, when open, queues the dropdown panel
// as a portal (see elements.RegisterPortal) anchored just below it, so
// it paints above everything else regardless of where Layout is
// called from.
func (s *Select) Layout(gtx layout.Context) layout.Dimensions {
	st := styles.Merge(styles.Styles{
		Background:   properties.CalcColorReverse(defaultBackground),
		Color:        properties.CalcColorReverse(defaultText),
		BorderWidth:  1,
		BorderColor:  properties.CalcColorReverse(defaultBorder),
		BorderRadius: 6,
		Padding:      properties.SymmetricInset(12, 8),
	}, s.st)

	if s.field.Clicked(gtx) {
		s.open = !s.open
	}

	// Field.
	fieldDims := s.field.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if s.st.Width > 0 {
			gtx.Constraints.Min.X = int(s.st.Width)
			gtx.Constraints.Max.X = int(s.st.Width)
		}
		fieldSt := st
		if s.field.Hovered() {
			fieldSt.BorderColor = kit.Hovered(st.BorderColor)
			pointer.CursorPointer.Add(gtx.Ops)
		}
		row := elements.Row(styles.Styles{Align: properties.AlignCenter},
			elements.Box(styles.Styles{FlexGrow: 1}, elements.Text(s.Value())),
			elements.Box(styles.Styles{Width: 4}),
			icon.Material("chevron-down", styles.Styles{Color: st.Color, Width: 16, Height: 16}),
		)
		return elements.Box(fieldSt, row).Layout(gtx)
	})
	// Keep the field clickable even inside a frameless window's
	// whole-window drag area - see button.Layout's identical comment.
	kit.BlockDrag(gtx, fieldDims.Size)

	if s.open && len(s.options) > 0 {
		fieldWidth := fieldDims.Size.X
		anchor := elements.AbsolutePosition().Add(image.Pt(0, fieldDims.Size.Y+4))
		elements.RegisterPortal(gtx, anchor, func() elements.Element {
			return elements.Raw(func(gtx layout.Context) layout.Dimensions {
				return s.dropdown(gtx, st, fieldWidth)
			})
		})
	}
	return fieldDims
}
