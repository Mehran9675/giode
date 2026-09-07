package properties

// Display selects the layout mode of a container.
type Display string

const (
	// Block flows children vertically, each filling the available
	// width.
	Block Display = "block"
	// Flex lays children out with the direction, align, justify and
	// gap properties.
	Flex Display = "flex"
	// None hides the element entirely: it takes no space and its
	// children are not laid out.
	None Display = "none"
)

// Mode resolves the display value, defaulting to block like CSS.
func (d Display) Mode() Display {
	if d == "" {
		return Block
	}
	return d
}
