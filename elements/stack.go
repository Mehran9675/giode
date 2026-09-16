package elements

import (
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// Stack lays children out vertically. It is Box with flex direction
// forced to column. Note that this is not Gio's layout.Stack.
func Stack(st styles.Styles, children ...Element) *BoxElement {
	st.FlexDirection = properties.FlexDirectionColumn
	return Box(st, children...)
}
