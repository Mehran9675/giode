package markdown

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/mehran9675/giode/elements"
	"github.com/mehran9675/giode/styles/properties"
)

// quote draws a block quote: content indented past a bar as tall as
// the content.
func quote(content elements.Element) elements.Element {
	return elements.Raw(func(gtx layout.Context) layout.Dimensions {
		bar, gap := gtx.Dp(3), gtx.Dp(10)
		cgtx := gtx
		cgtx.Constraints.Min.X = max(0, cgtx.Constraints.Min.X-bar-gap)
		cgtx.Constraints.Max.X = max(0, cgtx.Constraints.Max.X-bar-gap)
		off := image.Pt(bar+gap, 0)
		track := elements.TrackChild()
		macro := op.Record(gtx.Ops)
		dims := content.Layout(cgtx)
		call := macro.Stop()
		track.End()
		track.Place(off)

		func() {
			defer clip.UniformRRect(image.Rect(0, 0, bar, dims.Size.Y), bar/2).Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, properties.CalcColor("theme.border"))
		}()
		func() {
			defer op.Offset(off).Push(gtx.Ops).Pop()
			call.Add(gtx.Ops)
		}()
		return layout.Dimensions{Size: image.Pt(dims.Size.X+bar+gap, dims.Size.Y), Baseline: dims.Baseline}
	})
}

// codeScroll and tableScroll hold wide code and tables.
func codeScroll(el elements.Element) elements.Element  { return el }
func tableScroll(el elements.Element) elements.Element { return el }
