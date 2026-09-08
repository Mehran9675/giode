package router

import (
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

var defaultLinkColor = color.NRGBA{R: 0x60, G: 0xa5, B: 0xfa, A: 0xff}

// Link returns a clickable element that navigates to path when
// clicked.
func (r *Router) Link(path, text string, st ...styles.Styles) elements.Element {
	key := path + "\x00" + text
	click, ok := r.linkClicks[key]
	if !ok {
		click = new(widget.Clickable)
		r.linkClicks[key] = click
	}
	var s styles.Styles
	if len(st) > 0 {
		s = st[0]
	}
	return elements.Raw(func(gtx layout.Context) layout.Dimensions {
		if click.Clicked(gtx) {
			r.Navigate(path)
		}
		return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			col := properties.CalcColor(s.Color)
			if col.A == 0 {
				col = defaultLinkColor
			}
			if click.Hovered() {
				pointer.CursorPointer.Add(gtx.Ops)
			}
			return elements.Box(styles.Styles{
				Color:      properties.CalcColorReverse(col),
				FontSize:   s.FontSize,
				FontWeight: properties.FontWeightMedium,
			}, elements.Text(text)).Layout(gtx)
		})
	})
}
