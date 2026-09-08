package elements

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/styles/properties"
)

// layoutFlex lays children along an axis with align, justify and gap
// distribution, wrapping onto new lines when Wrap is active and the
// children overflow the main axis.
func (b *boxEl) layoutFlex(gtx layout.Context, children []Element) layout.Dimensions {
	st := b.st
	axis := st.FlexDirection.Axis()
	mainMax := mainConstraint(gtx.Constraints, axis)

	// Neutralize the main-axis minimum when the container hugs that
	// axis: Gio passes the parent's cross minimum through to rigid
	// children, which would otherwise inflate hug-sized containers
	// via the spacing distribution.
	hugMain := func() bool {
		if axis == layout.Horizontal {
			return st.Width == 0 && st.WidthPct == 0
		}
		return st.Height == 0 && st.HeightPct == 0
	}
	if hugMain() {
		if axis == layout.Horizontal {
			gtx.Constraints.Min.X = 0
		} else {
			gtx.Constraints.Min.Y = 0
		}
	}

	if st.Wrap.Resolved(axis) && mainMax < unbounded {
		return b.layoutFlexWrap(gtx, children, axis, mainMax)
	}

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
	dims := layout.Flex{
		Axis:      axis,
		Spacing:   st.Justify.FlexSpacing(),
		Alignment: st.Align.FlexAlignment(),
		Gap:       int(st.Gap),
	}.Layout(gtx, kids...)
	return dims
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

// wrapItem is one measured child: its recorded paint ops (so it is
// laid out exactly once) plus its size resolved into main/cross axis
// terms.
type wrapItem struct {
	ops         op.CallOp
	main, cross int
}

// wrapLine is one run of items that fit within the main axis.
type wrapLine struct {
	items     []wrapItem
	used      int // sum of item main sizes plus internal gaps
	thickness int // the line's extent on the cross axis
}

// layoutFlexWrap measures every child once (each is laid out exactly
// once, so stateful children behave correctly), partitions them into
// lines that fit mainMax, and stacks the lines on the cross axis.
//
// Percentage-sized children resolve to an exact main-axis size before
// measuring. FlexGrow has no effect here: a wrapped line's leftover
// space cannot be redistributed after a child has already been laid
// out once, so Justify only spaces items apart (as it already does
// for rigid children without wrap), it does not grow them.
func (b *boxEl) layoutFlexWrap(gtx layout.Context, children []Element, axis layout.Axis, mainMax int) layout.Dimensions {
	st := b.st
	gap := int(st.Gap)

	items := make([]wrapItem, 0, len(children))
	for _, child := range children {
		el := child
		if pct := flexPercentOf(child, axis); pct > 0 {
			el = &mainSized{el: child, axis: axis, size: mainMax * pct / 100}
		}
		cgtx := gtx
		if axis == layout.Horizontal {
			cgtx.Constraints.Min.X = 0
			cgtx.Constraints.Max.X = mainMax
		} else {
			cgtx.Constraints.Min.Y = 0
			cgtx.Constraints.Max.Y = mainMax
		}
		macro := op.Record(gtx.Ops)
		dims := el.Layout(cgtx)
		ops := macro.Stop()
		items = append(items, wrapItem{ops: ops, main: axisMain(dims.Size, axis), cross: axisCross(dims.Size, axis)})
	}

	var lines []wrapLine
	var cur wrapLine
	for _, it := range items {
		add := it.main
		if len(cur.items) > 0 {
			add += gap
		}
		if len(cur.items) > 0 && cur.used+add > mainMax {
			lines = append(lines, cur)
			cur = wrapLine{}
			add = it.main
		}
		cur.items = append(cur.items, it)
		cur.used += add
		if it.cross > cur.thickness {
			cur.thickness = it.cross
		}
	}
	if len(cur.items) > 0 {
		lines = append(lines, cur)
	}

	crossOffset := 0
	for li, ln := range lines {
		leftover := mainMax - ln.used
		lead, between := justifySpacing(st.Justify, leftover, len(ln.items))
		mainOffset := lead
		for i, it := range ln.items {
			crossPos := alignOffset(st.Align, ln.thickness, it.cross)
			off := op.Offset(axisPoint(axis, mainOffset, crossOffset+crossPos)).Push(gtx.Ops)
			it.ops.Add(gtx.Ops)
			off.Pop()
			mainOffset += it.main
			if i < len(ln.items)-1 {
				mainOffset += gap + between
			}
		}
		crossOffset += ln.thickness
		if li < len(lines)-1 {
			crossOffset += gap
		}
	}

	return layout.Dimensions{Size: axisPoint(axis, mainMax, crossOffset)}
}

// justifySpacing splits leftover main-axis space into a leading
// offset and extra per-gap spacing, mirroring Gio's Spacing modes.
func justifySpacing(justify properties.Justify, leftover, n int) (lead, between int) {
	if leftover <= 0 || n <= 0 {
		return 0, 0
	}
	switch justify {
	case properties.JustifyEnd:
		return leftover, 0
	case properties.JustifyCenter:
		return leftover / 2, 0
	case properties.JustifySpaceBetween:
		if n <= 1 {
			return 0, 0
		}
		return 0, leftover / (n - 1)
	case properties.JustifySpaceAround:
		g := leftover / n
		return g / 2, g
	case properties.JustifySpaceEvenly:
		g := leftover / (n + 1)
		return g, g
	default: // JustifyStart
		return 0, 0
	}
}

// alignOffset positions an item of the given cross-axis extent within
// a line of the given thickness.
func alignOffset(align properties.Align, thickness, cross int) int {
	switch align {
	case properties.AlignCenter:
		return (thickness - cross) / 2
	case properties.AlignEnd:
		return thickness - cross
	default:
		return 0
	}
}

// axisMain and axisCross read a size along and across axis.
func axisMain(p image.Point, axis layout.Axis) int {
	if axis == layout.Horizontal {
		return p.X
	}
	return p.Y
}

func axisCross(p image.Point, axis layout.Axis) int {
	if axis == layout.Horizontal {
		return p.Y
	}
	return p.X
}

// axisPoint builds a point from a main/cross pair for axis.
func axisPoint(axis layout.Axis, main, cross int) image.Point {
	if axis == layout.Horizontal {
		return image.Pt(main, cross)
	}
	return image.Pt(cross, main)
}
