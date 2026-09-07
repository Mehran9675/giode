package properties

// Position selects how an element is placed (CSS position).
type Position string

const (
	// PositionStatic is the default: the element takes part in the
	// normal flow.
	PositionStatic Position = "static"
	// PositionRelative behaves like static but can act as the
	// containing block of absolute descendants.
	PositionRelative Position = "relative"
	// PositionAbsolute removes the element from the flow and places
	// it at the Top/Right/Bottom/Left offsets of its parent box,
	// painting above the flow content.
	PositionAbsolute Position = "absolute"
)

// Mode resolves the value, defaulting to static.
func (p Position) Mode() Position {
	if p == "" {
		return PositionStatic
	}
	return p
}
