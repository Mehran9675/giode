package elements

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

// layoutAbsolute places absolutely positioned Box children at their
// Top/Left (or Right/Bottom) offsets within the parent's content box,
// painting above the flow content. The children do not affect the
// parent's size.
func (b *BoxElement) layoutAbsolute(gtx layout.Context, children []Element, size image.Point) {
	for _, child := range children {
		box, ok := child.(*BoxElement)
		if !ok {
			continue
		}
		st := box.st.Px(gtx.Metric)

		cgtx := gtx
		cgtx.Constraints.Min = image.Point{}
		cgtx.Constraints.Max = size
		if st.Left > 0 {
			cgtx.Constraints.Max.X -= int(st.Left)
		}
		if st.Right > 0 {
			cgtx.Constraints.Max.X -= int(st.Right)
		}
		if st.Top > 0 {
			cgtx.Constraints.Max.Y -= int(st.Top)
		}
		if st.Bottom > 0 {
			cgtx.Constraints.Max.Y -= int(st.Bottom)
		}
		if cgtx.Constraints.Max.X < 0 {
			cgtx.Constraints.Max.X = 0
		}
		if cgtx.Constraints.Max.Y < 0 {
			cgtx.Constraints.Max.Y = 0
		}

		// Left/top-anchored position is known upfront; right/bottom
		// anchoring depends on the child's own size, so it can only be
		// resolved after Layout returns. The absolute-position tracker
		// (see absolute.go) is pushed with the upfront position, and
		// what the child registered is moved by the difference once
		// the final one is known (see ChildTracker).
		x := int(st.Left)
		y := int(st.Top)
		upfront := image.Pt(x, y)

		popAbs := PushOffset(upfront)
		track := TrackChild()
		macro := op.Record(gtx.Ops)
		dims := box.Layout(cgtx)
		childOps := macro.Stop()
		track.End()
		popAbs()

		if st.Left == 0 && st.Right > 0 {
			x = size.X - int(st.Right) - dims.Size.X
		}
		if st.Top == 0 && st.Bottom > 0 {
			y = size.Y - int(st.Bottom) - dims.Size.Y
		}
		track.Place(image.Pt(x, y).Sub(upfront))

		off := op.Offset(image.Pt(x, y)).Push(gtx.Ops)
		childOps.Add(gtx.Ops)
		off.Pop()
	}
}
