package elements

import (
	"sort"

	"gioui.org/gesture"
	"gioui.org/layout"

	"github.com/mehran9675/giode/styles"
	"github.com/mehran9675/giode/styles/properties"
)

// BoxElement is a container element: a flex container and a
// containing block for absolutely positioned children.
type BoxElement struct {
	st       styles.Styles
	children []Element

	// Hover styles, set with OnHover. The hovered field caches the
	// last observed hover state.
	hover   *gesture.Hover
	hoverSt styles.Styles
	hovered bool
	hoverAnim
}

// Box is a container element: children are laid along FlexDirection
// with Align, Justify, Gap and Wrap. Chain OnHover to add hover
// styles.
func Box(st styles.Styles, children ...Element) *BoxElement {
	return &BoxElement{st: st, children: children}
}

// flexible marks elements that carry layout hints: how much they
// grow on the flex main axis, their percentage size and their
// z-index.
type flexible interface {
	flexWeight(axis layout.Axis) int
	percentOf(axis layout.Axis) int
	zIndex() int
}

func (b *BoxElement) flexWeight(axis layout.Axis) int {
	return weightOf(b.st, axis)
}

func (b *BoxElement) percentOf(axis layout.Axis) int {
	return percentOf(b.st, axis)
}

func (b *BoxElement) zIndex() int {
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
func (b *BoxElement) split() (flow, absolute []Element) {
	// Fast path, the overwhelmingly common case: no absolute children
	// and no z-index, so the children are used as they are without
	// allocating.
	simple := true
	for _, child := range b.children {
		if zIndexOf(child) != 0 {
			simple = false
			break
		}
		if box, ok := child.(*BoxElement); ok && box.st.Position.Mode() == properties.PositionAbsolute {
			simple = false
			break
		}
	}
	if simple {
		return b.children, nil
	}
	flow = make([]Element, 0, len(b.children))
	for _, child := range b.children {
		if box, ok := child.(*BoxElement); ok && box.st.Position.Mode() == properties.PositionAbsolute {
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

// minConstraint returns the main-axis minimum of c.
func minConstraint(c layout.Constraints, axis layout.Axis) int {
	if axis == layout.Horizontal {
		return c.Min.X
	}
	return c.Min.Y
}

// insetConstraints shrinks c by the inset on each side. Min shrinks by
// the same total inset as Max (clamped at 0), so an exact constraint
// (Min == Max, from a fixed or filled size) stays exact after the
// inset is removed instead of drifting.
func insetConstraints(c layout.Constraints, i properties.Inset) layout.Constraints {
	c.Min.X -= i.Left + i.Right
	c.Max.X -= i.Left + i.Right
	c.Min.Y -= i.Top + i.Bottom
	c.Max.Y -= i.Top + i.Bottom
	if c.Min.X < 0 {
		c.Min.X = 0
	}
	if c.Min.Y < 0 {
		c.Min.Y = 0
	}
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
