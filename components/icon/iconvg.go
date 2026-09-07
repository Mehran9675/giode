package icon

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// IconVG renders IconVG data, Gio's native vector icon format, sized
// to styles.Width.
func IconVG(data []byte, st ...styles.Styles) elements.Element {
	ic, err := widget.NewIcon(data)
	if err != nil {
		panic("giode: invalid iconvg: " + err.Error())
	}
	return &iconvgEl{icon: ic, st: first(st)}
}

// iconvgEl adapts widget.Icon to an element.
type iconvgEl struct {
	icon *widget.Icon
	st   styles.Styles
}

// Layout renders the icon at the given size.
func (i *iconvgEl) Layout(gtx layout.Context) layout.Dimensions {
	size := int(i.st.Width)
	if size <= 0 {
		size = 24
	}
	gtx.Constraints.Min.X = size
	gtx.Constraints.Max.X = size
	gtx.Constraints.Min.Y = size
	gtx.Constraints.Max.Y = size
	return i.icon.Layout(gtx, properties.ResolveColor(i.st.Color))
}
