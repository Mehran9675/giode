package elements

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/mehran9675/giode/styles"
)

var defaultDividerColor = color.NRGBA{R: 0x33, G: 0x41, B: 0x55, A: 0xff}

type dividerEl struct {
	st styles.Styles
}

// Divider is a thin horizontal line. Its color defaults to a subtle
// gray and its height to 1px; both can be overridden with styles.
func Divider(st ...styles.Styles) Element {
	var s styles.Styles
	if len(st) > 0 {
		s = st[0]
	}
	return &dividerEl{st: s}
}

func (d *dividerEl) Layout(gtx layout.Context) layout.Dimensions {
	height := int(d.st.Height)
	if height <= 0 {
		height = 1
	}
	width := gtx.Constraints.Max.X
	if width < gtx.Constraints.Min.X {
		width = gtx.Constraints.Min.X
	}
	size := image.Pt(width, height)
	c := d.st.Color
	if c.A == 0 {
		c = defaultDividerColor
	}
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, c)
	return layout.Dimensions{Size: size}
}
