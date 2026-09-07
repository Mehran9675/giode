package elements

import (
	"gioui.org/layout"
)

// layoutFlex lays children along an axis with align, justify and gap
// distribution.
func (b *boxEl) layoutFlex(gtx layout.Context, children []Element) layout.Dimensions {
	st := b.st
	axis := st.Direction.Axis()
	mainMax := mainConstraint(gtx.Constraints, axis)
	kids := make([]layout.FlexChild, 0, len(children))
	for _, child := range children {
		if pct := flexPercentOf(child, axis); pct > 0 && mainMax < unbounded {
			size := mainMax * pct / 100
			kids = append(kids, layout.Rigid((&mainSized{el: child, axis: axis, size: size}).Layout))
		} else if w := flexWeightOf(child, axis); w > 0 {
			kids = append(kids, layout.Flexed(float32(w), child.Layout))
		} else {
			kids = append(kids, layout.Rigid(child.Layout))
		}
	}
	return layout.Flex{
		Axis:      axis,
		Spacing:   st.Justify.FlexSpacing(),
		Alignment: st.Align.FlexAlignment(),
		Gap:       int(st.Gap),
	}.Layout(gtx, kids...)
}

// mainSized wraps an element to force its size on the flex main axis,
// used for percentage-sized flex children.
type mainSized struct {
	el   Element
	axis layout.Axis
	size int
}

func (s *mainSized) Layout(gtx layout.Context) layout.Dimensions {
	if s.axis == layout.Horizontal {
		gtx.Constraints.Min.X = s.size
		gtx.Constraints.Max.X = s.size
	} else {
		gtx.Constraints.Min.Y = s.size
		gtx.Constraints.Max.Y = s.size
	}
	return s.el.Layout(gtx)
}
