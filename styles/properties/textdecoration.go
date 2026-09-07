package properties

// TextDecoration decorates text (CSS text-decoration).
type TextDecoration string

const (
	// TextDecorationNone draws plain text.
	TextDecorationNone TextDecoration = "none"
	// TextDecorationUnderline draws a line under the text.
	TextDecorationUnderline TextDecoration = "underline"
	// TextDecorationLineThrough draws a line through the text.
	TextDecorationLineThrough TextDecoration = "line-through"
)
