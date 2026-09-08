package properties

import "gioui.org/layout"

// Wrap controls whether flex children continue onto a new line (CSS
// flex-wrap) when they no longer fit the container's main axis.
type Wrap string

const (
	// WrapNoWrap keeps every child on a single line, even if that
	// overflows the container.
	WrapNoWrap Wrap = "nowrap"
	// WrapWrap continues overflowing children onto a new line, stacked
	// on the cross axis.
	WrapWrap Wrap = "wrap"
)

// Resolved reports whether wrapping is active for a container with
// the given main axis. Unlike CSS (which defaults to nowrap), an
// unset Wrap defaults to wrapping on the horizontal axis (Row) and no
// wrapping on the vertical axis (Column); either can be overridden
// explicitly.
func (w Wrap) Resolved(axis layout.Axis) bool {
	switch w {
	case WrapWrap:
		return true
	case WrapNoWrap:
		return false
	default:
		return axis == layout.Horizontal
	}
}
