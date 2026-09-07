package tabs

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/mehran9675/giode/components/kit"
	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// Layout lays out the tab bar and the active view below it.
func (t *Tabs) Layout(gtx layout.Context) layout.Dimensions {
	if len(t.tabs) == 0 {
		return layout.Dimensions{}
	}
	active := t.st.Color
	if active.A == 0 {
		active = defaultActive
	}
	activeText := t.st.BorderColor
	if activeText.A == 0 {
		activeText = defaultActiveText
	}

	var size image.Point
	dims := elements.Flex(styles.Styles{Gap: 4}, t.tabButtons(active, activeText)...).Layout(gtx)
	size.Y += dims.Size.Y

	if t.selected < len(t.tabs) {
		gtx.Constraints.Min.Y -= dims.Size.Y
		gtx.Constraints.Max.Y -= dims.Size.Y
		off := op.Offset(image.Pt(0, dims.Size.Y)).Push(gtx.Ops)
		contentDims := t.tabs[t.selected].View().Layout(gtx)
		off.Pop()
		size.Y += contentDims.Size.Y
		if contentDims.Size.X > size.X {
			size.X = contentDims.Size.X
		}
	}
	return layout.Dimensions{Size: size}
}

// tabButtons returns one header element per tab.
func (t *Tabs) tabButtons(active, activeText color.NRGBA) []elements.Element {
	els := make([]elements.Element, 0, len(t.tabs))
	for i := range t.tabs {
		els = append(els, t.tabButton(i, active, activeText))
	}
	return els
}

// tabButton renders one clickable tab header.
func (t *Tabs) tabButton(i int, active, activeText color.NRGBA) elements.Element {
	click := t.clicks[i]
	return elements.Raw(func(gtx layout.Context) layout.Dimensions {
		if click.Clicked(gtx) {
			t.selected = i
			if t.onChange != nil {
				t.onChange(i)
			}
		}
		return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			sel := i == t.selected
			textColor := defaultText
			if sel {
				textColor = activeText
			} else if click.Hovered() {
				textColor = kit.Hovered(defaultText)
			}
			if click.Hovered() {
				pointer.CursorPointer.Add(gtx.Ops)
			}
			st := styles.Styles{
				Padding: properties.SymmetricInset(12, 8),
			}
			content := elements.Text(t.tabs[i].Label, styles.Styles{Color: textColor, FontSize: 14})
			if !sel {
				return elements.Box(st, content).Layout(gtx)
			}
			// Selected tab: label + 2px underline.
			macro := op.Record(gtx.Ops)
			dims := elements.Box(st, content).Layout(gtx)
			contentOps := macro.Stop()
			defer clip.Rect{Max: dims.Size}.Push(gtx.Ops).Pop()
			contentOps.Add(gtx.Ops)
			off := op.Offset(image.Pt(0, dims.Size.Y-2)).Push(gtx.Ops)
			paint.FillShape(gtx.Ops, active, clip.Rect{Max: image.Pt(dims.Size.X, 2)}.Op())
			off.Pop()
			return dims
		})
	})
}
