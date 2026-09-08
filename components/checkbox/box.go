package checkbox

import (
	"image"

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
	defaultBox     = "#475569"
	defaultChecked = "#3b82f6"
	defaultText    = "#e2e8f0"
	defaultIcon    = "#ffffff"
)

const (
	defaultBoxSize     = 18
	defaultBorderWidth = 2
	defaultRadius      = 4
)

// box renders the square indicator: filled with an icon when checked,
// outlined when not.
func (c *Checkbox) box() elements.Element {
	return elements.Raw(func(gtx layout.Context) layout.Dimensions {
		size := c.boxSize()
		gtx.Constraints.Min = size
		gtx.Constraints.Max = size

		radius := int(c.st.BorderRadius)
		if radius == 0 {
			radius = defaultRadius
		}
		borderWidth := c.st.BorderWidth
		if borderWidth == 0 {
			borderWidth = defaultBorderWidth
		}

		checked := c.bool.Value
		fill := stOrDefault(c.st.Background, defaultChecked)
		border := stOrDefault(c.st.BorderColor, defaultBox)
		fillColor := properties.CalcColor(fill)
		borderColor := properties.CalcColor(border)
		switch {
		case !checked && c.bool.Hovered():
			borderColor = kit.Hovered(borderColor)
		case checked && c.bool.Hovered():
			fillColor = kit.Hovered(fillColor)
		}

		defer clip.UniformRRect(image.Rectangle{Max: size}, radius).Push(gtx.Ops).Pop()
		if checked {
			paint.Fill(gtx.Ops, fillColor)
			iconSt := styles.Merge(styles.Styles{
				Color:  defaultIcon,
				Width:  properties.Width(size.X),
				Height: properties.Height(size.Y),
			}, c.iconSt)
			icon.Material(c.iconName, iconSt).Layout(gtx)
		} else {
			properties.PaintBorder(gtx.Ops, size, properties.BorderWidth(borderWidth), properties.CalcColorReverse(borderColor), properties.BorderRadius(radius))
		}
		return layout.Dimensions{Size: size}
	})
}

// boxSize resolves the indicator size from Width/Height, defaulting
// to a square 18x18 box.
func (c *Checkbox) boxSize() image.Point {
	w := int(c.st.Width)
	h := int(c.st.Height)
	if w <= 0 {
		w = defaultBoxSize
	}
	if h <= 0 {
		h = defaultBoxSize
	}
	return image.Pt(w, h)
}

// stOrDefault returns c unless it is unset, in which case it returns
// def.
func stOrDefault(c, def properties.Color) properties.Color {
	if c == "" {
		return def
	}
	return c
}
