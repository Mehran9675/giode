package button

import (
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"

	"github.com/mehran9675/giode/components/internal/colorutil"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

var defaultBackground = color.NRGBA{R: 0x4b, G: 0x55, B: 0x63, A: 0xff}

// layoutVisual renders the button face inside the clickable region.
func (b *Button) layoutVisual(gtx layout.Context) layout.Dimensions {
	st := b.st
	if st.Padding == (properties.Inset{}) {
		st.Padding = properties.SymmetricInset(16, 8)
	}
	bg := st.Background
	if bg.A == 0 {
		bg = defaultBackground
	}
	switch {
	case b.clickable.Pressed():
		bg = colorutil.Pressed(bg)
	case b.clickable.Hovered():
		bg = colorutil.Hovered(bg)
		pointer.CursorPointer.Add(gtx.Ops)
	}
	st.Background = bg
	label := elements.Text(b.label, styles.Styles{
		Color:    st.Color,
		FontSize: st.FontSize,
	})
	return elements.Box(st, label).Layout(gtx)
}
