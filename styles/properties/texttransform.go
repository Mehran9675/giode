package properties

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// TextTransform changes the capitalization of text (CSS
// text-transform).
type TextTransform string

const (
	TextTransformNone       TextTransform = "none"
	TextTransformUppercase  TextTransform = "uppercase"
	TextTransformLowercase  TextTransform = "lowercase"
	TextTransformCapitalize TextTransform = "capitalize"
)

// Apply returns s transformed.
func (t TextTransform) Apply(s string) string {
	switch t {
	case TextTransformUppercase:
		return strings.ToUpper(s)
	case TextTransformLowercase:
		return strings.ToLower(s)
	case TextTransformCapitalize:
		var b strings.Builder
		b.Grow(len(s))
		start := true
		for len(s) > 0 {
			r, n := utf8.DecodeRuneInString(s)
			s = s[n:]
			if start && unicode.IsLetter(r) {
				r = unicode.ToTitle(r)
			}
			start = unicode.IsSpace(r) || r == '-'
			b.WriteRune(r)
		}
		return b.String()
	default:
		return s
	}
}

// TextOverflow selects what happens to text cut off by MaxLines (CSS
// text-overflow).
type TextOverflow string

const (
	// TextOverflowClip cuts the text without a marker (the default).
	TextOverflowClip TextOverflow = "clip"
	// TextOverflowEllipsis ends truncated text with "…". Without an
	// explicit MaxLines it implies a single line.
	TextOverflowEllipsis TextOverflow = "ellipsis"
)
