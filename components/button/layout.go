package button

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"

	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

var defaultBackground = "#000000"

// layoutVisual renders the button face inside the clickable region.
func (b *Button) layoutVisual(gtx layout.Context) layout.Dimensions {
	st := styles.Merge(styles.Styles{
		Padding: properties.SymmetricInset(16, 8),
		Justify: properties.JustifyCenter,
		Align:   properties.AlignCenter,
	}, b.st)
	bg := properties.CalcColor(st.Background)
	if bg.A == 0 {
		bg = properties.CalcColor(defaultBackground)
	}
	switch {
	case b.clickable.Pressed():
		bg = kit.Pressed(bg)
	case b.clickable.Hovered():
		bg = kit.Hovered(bg)
		pointer.CursorPointer.Add(gtx.Ops)
	}
	st.Background = properties.CalcColorReverse(bg)
	label := elements.Text(b.label, styles.Styles{
		Color:     st.Color,
		FontSize:  st.FontSize,
		TextAlign: st.TextAlign,
	})
	return elements.Box(st, label).Layout(gtx)
}
