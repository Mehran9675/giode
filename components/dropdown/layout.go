package dropdown

import (
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/components/icon"
	"github.com/mehran9675/giode/components/internal/overlay"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// Layout lays out the field and, when open, the dropdown below it.
func (s *Select) Layout(gtx layout.Context) layout.Dimensions {
	st := styles.Merge(styles.Styles{
		Background:   defaultBackground,
		Color:        defaultText,
		BorderWidth:  1,
		BorderColor:  defaultBorder,
		BorderRadius: 6,
		Padding:      properties.SymmetricInset(12, 8),
	}, s.st)

	if s.field.Clicked(gtx) {
		s.open = !s.open
	}

	// Field.
	macro := op.Record(gtx.Ops)
	fieldDims := s.field.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if s.st.Width > 0 {
			gtx.Constraints.Min.X = int(s.st.Width)
			gtx.Constraints.Max.X = int(s.st.Width)
		}
		if s.field.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
		}
		row := elements.Flex(styles.Styles{Align: properties.AlignCenter},
			elements.Box(styles.Styles{FlexGrow: 1}, elements.Text(s.Value())),
			elements.Box(styles.Styles{Width: 4}),
			icon.Material("chevron-down", styles.Styles{Color: st.Color, Width: 16, Height: 16}),
		)
		return elements.Box(st, row).Layout(gtx)
	})
	fieldOps := macro.Stop()

	// Dropdown.
	var dropdownDims layout.Dimensions
	var dropdownOps op.CallOp
	if s.open && len(s.options) > 0 {
		m := op.Record(gtx.Ops)
		dropdownDims = s.dropdown(gtx, st, fieldDims.Size.X)
		dropdownOps = m.Stop()
	}

	fieldOps.Add(gtx.Ops)
	if s.open && len(s.options) > 0 {
		panelSize := f32.Pt(float32(dropdownDims.Size.X), float32(dropdownDims.Size.Y))
		pos := overlay.ClampPos(
			f32.Pt(float32(gtx.Constraints.Max.X), float32(gtx.Constraints.Max.Y)),
			panelSize,
			f32.Pt(0, float32(fieldDims.Size.Y)+4),
			overlay.Margin,
		)
		off := op.Offset(pos.Round()).Push(gtx.Ops)
		dropdownOps.Add(gtx.Ops)
		off.Pop()
	}
	return fieldDims
}
