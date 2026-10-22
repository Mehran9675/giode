package properties

// Display controls whether an element takes part in layout (CSS
// display).
type Display string

const (
	// DisplayDefault lays the element out normally.
	DisplayDefault Display = ""
	// DisplayNone removes the element: it takes no space and paints
	// nothing.
	DisplayNone Display = "none"
)

// Visibility controls whether an element paints while still taking
// part in layout (CSS visibility).
type Visibility string

const (
	// VisibilityVisible paints the element (the default).
	VisibilityVisible Visibility = "visible"
	// VisibilityHidden keeps the element's space but paints nothing
	// and ignores pointer input.
	VisibilityHidden Visibility = "hidden"
)
