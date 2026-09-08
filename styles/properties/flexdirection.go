package properties

import "gioui.org/layout"

// FlexDirection is the main axis of a flex container (CSS
// flex-direction).
type FlexDirection string

const (
	// FlexDirectionRow lays children left to right.
	FlexDirectionRow FlexDirection = "row"
	// FlexDirectionColumn lays children top to bottom.
	FlexDirectionColumn FlexDirection = "column"
)

// Axis converts the direction to a layout axis, defaulting to row.
func (d FlexDirection) Axis() layout.Axis {
	if d == FlexDirectionColumn {
		return layout.Vertical
	}
	return layout.Horizontal
}
