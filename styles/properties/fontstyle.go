package properties

import "gioui.org/font"

// FontStyle selects the text style (CSS font-style).
type FontStyle string

const (
	FontStyleNormal FontStyle = "normal"
	FontStyleItalic FontStyle = "italic"
)

// Style maps the value to a font style, defaulting to regular.
func (s FontStyle) Style() font.Style {
	if s == FontStyleItalic {
		return font.Italic
	}
	return font.Regular
}
