package checkbox

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/mehran9675/giode/components/icon"
	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

var (
	defaultBox     = color.NRGBA{R: 0x47, G: 0x55, B: 0x69, A: 0xff}
	defaultChecked = color.NRGBA{R: 0x3b, G: 0x82, B: 0xf6, A: 0xff}
	defaultText    = color.NRGBA{R: 0xe2, G: 0xe8, B: 0xf0, A: 0xff}
)

const boxSize = 18

// box renders the square indicator.
func (c *Checkbox) box() elements.Element {
	return elements.Raw(func(gtx layout.Context) layout.Dimensions {
		size := image.Pt(boxSize, boxSize)
		gtx.Constraints.Min = size
		gtx.Constraints.Max = size
		checked := c.bool.Value
		fill := stOrDefault(c.st.Background, defaultChecked)
		border := stOrDefault(c.st.BorderColor, defaultBox)
		switch {
		case !checked && c.bool.Hovered():
			border = kit.Hovered(border)
		case checked && c.bool.Hovered():
			fill = kit.Hovered(fill)
		}
		if !checked {
			fill = color.NRGBA{}
		}
		defer clip.UniformRRect(image.Rectangle{Max: size}, 4).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, fill)
		if !checked {
			properties.PaintBorder(gtx.Ops, size, 2, border, 4)
		} else {
			icon.Material("check", styles.Styles{Color: color.NRGBA{A: 0xff}, Width: boxSize, Height: boxSize}).Layout(gtx)
		}
		return layout.Dimensions{Size: size}
	})
}

func stOrDefault(c, def color.NRGBA) color.NRGBA {
	if c.A == 0 {
		return def
	}
	return c
}
