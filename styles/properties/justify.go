package properties

import "gioui.org/layout"

// Justify distributes children on the main axis of a flex container
// (CSS justify-content).
type Justify string

const (
	JustifyStart        Justify = "start"
	JustifyCenter       Justify = "center"
	JustifyEnd          Justify = "end"
	JustifySpaceBetween Justify = "space-between"
	JustifySpaceAround  Justify = "space-around"
	JustifySpaceEvenly  Justify = "space-evenly"
)

// FlexSpacing maps the value to layout.Flex spacing. Note that like
// CSS, distribution only happens when the container has spare space on
// the main axis.
func (j Justify) FlexSpacing() layout.Spacing {
	switch j {
	case JustifyEnd:
		return layout.SpaceStart
	case JustifyCenter:
		return layout.SpaceSides
	case JustifySpaceBetween:
		return layout.SpaceBetween
	case JustifySpaceAround:
		return layout.SpaceAround
	case JustifySpaceEvenly:
		return layout.SpaceEvenly
	default:
		return layout.SpaceEnd
	}
}
