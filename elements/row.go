package elements

import (
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// Row lays children out horizontally. It is Box with flex direction
// forced to row.
func Row(st styles.Styles, children ...Element) *BoxElement {
	st.FlexDirection = properties.FlexDirectionRow
	return Box(st, children...)
}
