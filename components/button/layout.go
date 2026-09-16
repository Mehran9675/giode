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

	accent := st.Background
	if accent == "" {
		accent = defaultBackground
	}
	hovered := b.clickable.Hovered()
	pressed := b.clickable.Pressed()

	switch b.variant {
	case VariantOutlined:
		border := st.BorderColor
		if border == "" {
			border = accent
		}
		st.BorderColor = border
		if st.BorderWidth == 0 {
			st.BorderWidth = 1
		}
		st.Background = ""
		switch {
		case pressed:
			st.Background = kit.Alpha(border, 0.24)
		case hovered:
			st.Background = kit.Alpha(border, 0.12)
			pointer.CursorPointer.Add(gtx.Ops)
		}
	case VariantGhost:
		// Fully transparent at rest: only the label is visible until
		// hover reveals the box.
		st.BorderWidth = 0
		st.Background = ""
		switch {
		case pressed:
			st.Background = kit.Alpha(accent, 0.24)
		case hovered:
			st.Background = kit.Alpha(accent, 0.12)
			pointer.CursorPointer.Add(gtx.Ops)
		}
	default: // VariantContained, and the zero value for compatibility.
		switch {
		case pressed:
			accent = kit.Pressed(accent)
		case hovered:
			accent = kit.Hovered(accent)
			pointer.CursorPointer.Add(gtx.Ops)
		}
		st.Background = accent
	}

	label := elements.Text(b.label, styles.Styles{
		Color:     st.Color,
		FontSize:  st.FontSize,
		TextAlign: st.TextAlign,
	})
	return elements.Box(st, label).Layout(gtx)
}
