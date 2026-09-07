package elements

import (
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// Flex lays children out in a row or column with align, justify and
// gap distribution. It is Box with display forced to flex.
func Flex(st styles.Styles, children ...Element) Element {
	st.Display = properties.Flex
	return Box(st, children...)
}
