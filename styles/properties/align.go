package properties

import "gioui.org/layout"

// Align positions children on the cross axis of a flex container
// (CSS align-items).
type Align string

const (
	AlignStart  Align = "start"
	AlignCenter Align = "center"
	AlignEnd    Align = "end"
)

// FlexAlignment maps the value to a cross-axis layout alignment. The
// CSS default of "stretch" is not supported yet and behaves like
// "start".
func (a Align) FlexAlignment() layout.Alignment {
	switch a {
	case AlignCenter:
		return layout.Middle
	case AlignEnd:
		return layout.End
	default:
		return layout.Start
	}
}
