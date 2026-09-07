package properties

import "gioui.org/layout"

// Direction is the main axis of a flex container.
type Direction string

const (
	// Row lays children left to right.
	Row Direction = "row"
	// Column lays children top to bottom.
	Column Direction = "column"
)

// Axis converts the direction to a layout axis, defaulting to row.
func (d Direction) Axis() layout.Axis {
	if d == Column {
		return layout.Vertical
	}
	return layout.Horizontal
}
