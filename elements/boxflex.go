package elements

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"

	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// layoutFlex lays children along an axis with align, justify and gap
// distribution, wrapping onto new lines when Wrap is active and the
// children overflow the main axis. st may differ from b.st when hover
// styles are active.
func (b *BoxElement) layoutFlex(gtx layout.Context, st styles.Styles, children []Element) layout.Dimensions {
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
		return b.layoutFlexWrap(gtx, st, children, axis, mainMax)
	}

	// Gio passes the container's cross-axis minimum on to every child,
	// stretching it. That matches CSS's default (align-items: stretch)
	// for start alignment, but centered or end-aligned children must
	// keep their own size or there is nothing to align.
	unstretch := st.Align == properties.AlignCenter || st.Align == properties.AlignEnd
	var track *semTracker
	if semantics != nil {
		track = &semTracker{spans: make([]semSpan, len(children))}
	}
	kids := make([]layout.FlexChild, 0, len(children))
	for i, child := range children {
		w := child.Layout
		if pct := flexPercentOf(child, axis); pct > 0 && mainMax < unbounded {
			w = (&mainSized{el: child, axis: axis, size: mainMax * pct / 100}).Layout
		}
		if unstretch {
			w = unstretchCross(w, axis)
		}
		if track != nil {
			w = track.wrap(i, w)
		}
		if wt := flexWeightOf(child, axis); wt > 0 && flexPercentOf(child, axis) == 0 {
			kids = append(kids, layout.Flexed(float32(wt), w))
		} else {
			kids = append(kids, layout.Rigid(w))
		}
	}
	flex := layout.Flex{
		Axis:      axis,
		Spacing:   st.Justify.FlexSpacing(),
		Alignment: st.Align.FlexAlignment(),
		Gap:       int(st.Gap),
	}
	dims := flex.Layout(gtx, kids...)
	if track != nil {
		track.place(gtx.Constraints, flex)
	}
	return dims
}

// unstretchCross lets a child size itself on the cross axis.
func unstretchCross(w layout.Widget, axis layout.Axis) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if axis == layout.Horizontal {
			gtx.Constraints.Min.Y = 0
		} else {
			gtx.Constraints.Min.X = 0
		}
		return w(gtx)
	}
}

// semSpan is the range of semantic nodes one flex child recorded, and
// its measured size.
type semSpan struct {
	start, end int
	dims       layout.Dimensions
}

// semTracker fixes up the bounds of semantic nodes recorded by flex
// children. A child lays out before the flex container knows where it
// goes, so its nodes are recorded at the container's origin; once the
// container is done, each child's nodes are shifted by the offset Gio's
// flex gave it.
type semTracker struct {
	spans []semSpan
}

func (t *semTracker) wrap(i int, w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		t.spans[i].start = len(*semantics)
		dims := w(gtx)
		t.spans[i].end = len(*semantics)
		t.spans[i].dims = dims
		return dims
	}
}

// place replicates the placement step of layout.Flex.Layout.
func (t *semTracker) place(cs layout.Constraints, f layout.Flex) {
	if semantics == nil {
		return
	}
	n := len(t.spans)
	mainMin := axisMain(cs.Min, f.Axis)
	crossMin := axisCross(cs.Min, f.Axis)
	size, maxCross, maxBaseline := 0, crossMin, 0
	for _, s := range t.spans {
		size += axisMain(s.dims.Size, f.Axis)
		maxCross = max(maxCross, axisCross(s.dims.Size, f.Axis))
		maxBaseline = max(maxBaseline, s.dims.Size.Y-s.dims.Baseline)
	}
	if n > 1 && f.Gap > 0 {
		size += f.Gap * (n - 1)
	}
	space := max(0, mainMin-size)
	main := 0
	switch f.Spacing {
	case layout.SpaceSides:
		main += space / 2
	case layout.SpaceStart:
		main += space
	case layout.SpaceEvenly:
		main += space / (1 + n)
	case layout.SpaceAround:
		if n > 0 {
			main += space / (n * 2)
		}
	}
	nodes := *semantics
	for i, s := range t.spans {
		cross := 0
		switch f.Alignment {
		case layout.End:
			cross = maxCross - axisCross(s.dims.Size, f.Axis)
		case layout.Middle:
			cross = (maxCross - axisCross(s.dims.Size, f.Axis)) / 2
		case layout.Baseline:
			if f.Axis == layout.Horizontal {
				cross = maxBaseline - (s.dims.Size.Y - s.dims.Baseline)
			}
		}
		off := axisPoint(f.Axis, main, cross)
		for j := s.start; j < s.end && j < len(nodes); j++ {
			nodes[j].Bounds = nodes[j].Bounds.Add(off)
		}
		main += axisMain(s.dims.Size, f.Axis)
		if i < n-1 {
			main += f.Gap
			switch f.Spacing {
			case layout.SpaceEvenly:
				main += space / (1 + n)
			case layout.SpaceAround:
				main += space / n
			case layout.SpaceBetween:
				if n > 1 {
					main += space / (n - 1)
				}
			}
		}
	}
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
	// semStart and semEnd delimit the semantic nodes the item
	// recorded, shifted once its position is known.
	semStart, semEnd int
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
func (b *BoxElement) layoutFlexWrap(gtx layout.Context, st styles.Styles, children []Element, axis layout.Axis, mainMax int) layout.Dimensions {
	gap := int(st.Gap)

	// A container that hugs the main axis has no leftover space to
	// distribute: mainMax is only the space it may grow into, while
	// the container (and any clip around it, e.g. a Clickable's) ends
	// up as wide as its content. Justify must not push the content
	// across that unused space.
	hugMain := false
	if axis == layout.Horizontal {
		hugMain = st.Width == 0 && st.WidthPct == 0
	} else {
		hugMain = st.Height == 0 && st.HeightPct == 0
	}

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
		if st.Align == properties.AlignCenter || st.Align == properties.AlignEnd {
			// See layoutFlex: aligned items keep their own cross size.
			if axis == layout.Horizontal {
				cgtx.Constraints.Min.Y = 0
			} else {
				cgtx.Constraints.Min.X = 0
			}
		}
		semStart := 0
		if semantics != nil {
			semStart = len(*semantics)
		}
		macro := op.Record(gtx.Ops)
		dims := el.Layout(cgtx)
		ops := macro.Stop()
		it := wrapItem{ops: ops, main: axisMain(dims.Size, axis), cross: axisCross(dims.Size, axis)}
		if semantics != nil {
			it.semStart, it.semEnd = semStart, len(*semantics)
		}
		items = append(items, it)
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

	// Lines share any spare cross-axis space of a container with a
	// minimum cross size (a fixed Height on a Row, say), like CSS's
	// default align-content: stretch, so aligned items align within
	// the container rather than within their own line.
	if len(lines) > 0 {
		total := gap * (len(lines) - 1)
		for _, ln := range lines {
			total += ln.thickness
		}
		if extra := axisCross(gtx.Constraints.Min, axis) - total; extra > 0 {
			for i := range lines {
				lines[i].thickness += extra / len(lines)
			}
		}
	}

	crossOffset := 0
	contentMain := 0
	for li, ln := range lines {
		if ln.used > contentMain {
			contentMain = ln.used
		}
		leftover := mainMax - ln.used
		if hugMain {
			leftover = 0
		}
		lead, between := justifySpacing(st.Justify, leftover, len(ln.items))
		mainOffset := lead
		for i, it := range ln.items {
			crossPos := alignOffset(st.Align, ln.thickness, it.cross)
			pos := axisPoint(axis, mainOffset, crossOffset+crossPos)
			if semantics != nil {
				for j := it.semStart; j < it.semEnd; j++ {
					(*semantics)[j].Bounds = (*semantics)[j].Bounds.Add(pos)
				}
			}
			off := op.Offset(pos).Push(gtx.Ops)
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

	// Report the widest line's actual extent, not the full mainMax: a
	// hugging container (Min < mainMax, e.g. a button sized to its
	// label) must still hug its wrapped content instead of claiming
	// all the space it was merely allowed to wrap within. A container
	// required to fill or that's exactly sized (Min == mainMax) still
	// reports mainMax, since contentMain is clamped up to it.
	mainSize := contentMain
	if minMain := minConstraint(gtx.Constraints, axis); mainSize < minMain {
		mainSize = minMain
	}
	if mainSize > mainMax {
		mainSize = mainMax
	}

	// The cross extent is clamped the same way, so a container with a
	// forced cross size (e.g. Height: -1) still fills it even when its
	// wrapped content is shorter.
	crossSize := crossOffset
	if minCross := axisCross(gtx.Constraints.Min, axis); crossSize < minCross {
		crossSize = minCross
	}
	if maxCross := axisCross(gtx.Constraints.Max, axis); crossSize > maxCross {
		crossSize = maxCross
	}

	return layout.Dimensions{Size: axisPoint(axis, mainSize, crossSize)}
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
