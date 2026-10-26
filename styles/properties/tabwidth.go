package properties

import "strings"

// TabWidth is the distance between tab stops in text, in columns (CSS
// tab-size). Zero means 4. Fonts rarely have a glyph for the tab
// character, so tabs are expanded to spaces up to the next stop; in a
// monospaced font that lines code up exactly.
type TabWidth int

// Expand replaces each tab in s with spaces up to the next tab stop,
// counting columns from the start of each line.
func (w TabWidth) Expand(s string) string {
	if !strings.ContainsRune(s, '\t') {
		return s
	}
	n := int(w)
	if n <= 0 {
		n = 4
	}
	var b strings.Builder
	b.Grow(len(s) + 8)
	col := 0
	for _, r := range s {
		switch r {
		case '\t':
			pad := n - col%n
			b.WriteString(strings.Repeat(" ", pad))
			col += pad
		case '\n':
			b.WriteRune(r)
			col = 0
		default:
			b.WriteRune(r)
			col++
		}
	}
	return b.String()
}
