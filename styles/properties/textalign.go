package properties

import "gioui.org/text"

// TextAlign aligns text within its box (CSS text-align).
type TextAlign string

const (
	TextAlignStart  TextAlign = "start"
	TextAlignCenter TextAlign = "center"
	TextAlignEnd    TextAlign = "end"
)

// LayoutAlignment maps the value to a text alignment, defaulting to
// start.
func (a TextAlign) LayoutAlignment() text.Alignment {
	switch a {
	case TextAlignCenter:
		return text.Middle
	case TextAlignEnd:
		return text.End
	default:
		return text.Start
	}
}
