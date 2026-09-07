package elements

import (
	"sort"

	"gioui.org/layout"

	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// boxEl is a container element: block flow, flex layout or a
// containing block for absolutely positioned children.
type boxEl struct {
	st       styles.Styles
	children []Element
}

// Box is a container element. Its display property selects between
// block flow (children stacked vertically, filling the width), flex
// layout, and none (not laid out).
func Box(st styles.Styles, children ...Element) Element {
	return &boxEl{st: st, children: children}
}

// flexible marks elements that carry layout hints: how much they
// grow on the flex main axis, their percentage size and their
// z-index.
type flexible interface {
	flexWeight(axis layout.Axis) int
	percentOf(axis layout.Axis) int
	zIndex() int
}

func (b *boxEl) flexWeight(axis layout.Axis) int {
	return weightOf(b.st, axis)
}

func (b *boxEl) percentOf(axis layout.Axis) int {
	return percentOf(b.st, axis)
}

func (b *boxEl) zIndex() int {
	return int(b.st.ZIndex)
}

func weightOf(st styles.Styles, axis layout.Axis) int {
	if st.FlexGrow > 0 {
		return int(st.FlexGrow)
	}
	if axis == layout.Horizontal && st.Width == -1 {
		return 1
	}
	if axis == layout.Vertical && st.Height == -1 {
		return 1
	}
	return 0
}

func percentOf(st styles.Styles, axis layout.Axis) int {
	if axis == layout.Horizontal {
		return int(st.WidthPct)
	}
	return int(st.HeightPct)
}

func flexWeightOf(el Element, axis layout.Axis) int {
	if f, ok := el.(flexible); ok {
		return f.flexWeight(axis)
	}
	return 0
}

func flexPercentOf(el Element, axis layout.Axis) int {
	if f, ok := el.(flexible); ok {
		return f.percentOf(axis)
	}
	return 0
}

func zIndexOf(el Element) int {
	if f, ok := el.(flexible); ok {
		return f.zIndex()
	}
	return 0
}

// split partitions children into flow and absolutely positioned
// children, each sorted by z-index. Absolute positioning applies to
// Box children.
func (b *boxEl) split() (flow, absolute []Element) {
	flow = make([]Element, 0, len(b.children))
	for _, child := range b.children {
		if box, ok := child.(*boxEl); ok && box.st.Position.Mode() == properties.PositionAbsolute {
			absolute = append(absolute, child)
		} else {
			flow = append(flow, child)
		}
	}
	sortZ(flow)
	sortZ(absolute)
	return flow, absolute
}

func sortZ(els []Element) {
	sort.SliceStable(els, func(i, j int) bool {
		return zIndexOf(els[i]) < zIndexOf(els[j])
	})
}

// mainConstraint returns the main-axis maximum of c.
func mainConstraint(c layout.Constraints, axis layout.Axis) int {
	if axis == layout.Horizontal {
		return c.Max.X
	}
	return c.Max.Y
}

// insetConstraints shrinks c by the inset on each side.
func insetConstraints(c layout.Constraints, i properties.Inset) layout.Constraints {
	c.Min.X += i.Left
	c.Max.X -= i.Left + i.Right
	c.Min.Y += i.Top
	c.Max.Y -= i.Top + i.Bottom
	if c.Max.X < c.Min.X {
		c.Max.X = c.Min.X
	}
	if c.Max.Y < c.Min.Y {
		c.Max.Y = c.Min.Y
	}
	return c
}

// unbounded mirrors the sentinel Gio uses for infinite constraints.
const unbounded = 1e6
