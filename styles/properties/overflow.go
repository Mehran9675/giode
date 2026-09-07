package properties

// Overflow controls what happens to content outside an element's
// box.
type Overflow string

const (
	// OverflowVisible (default) lets content paint outside the box.
	OverflowVisible Overflow = "visible"
	// OverflowHidden clips content to the box.
	OverflowHidden Overflow = "hidden"
)

// Mode resolves the value, defaulting to visible.
func (o Overflow) Mode() Overflow {
	if o == "" {
		return OverflowVisible
	}
	return o
}
