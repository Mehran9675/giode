package elements

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

// layoutBlock stacks children vertically, each filling the available
// width like CSS block layout.
func (b *boxEl) layoutBlock(gtx layout.Context, children []Element) layout.Dimensions {
	if len(children) == 0 {
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	var size image.Point
	var offsets []op.TransformStack
	for i, child := range children {
		dims := child.Layout(gtx)
		size.X = max(size.X, dims.Size.X)
		size.Y += dims.Size.Y
		if i < len(children)-1 {
			gtx.Constraints.Min.Y -= dims.Size.Y
			gtx.Constraints.Max.Y -= dims.Size.Y
			offsets = append(offsets, op.Offset(image.Pt(0, dims.Size.Y)).Push(gtx.Ops))
		}
	}
	for i := len(offsets) - 1; i >= 0; i-- {
		offsets[i].Pop()
	}
	return layout.Dimensions{Size: size}
}
