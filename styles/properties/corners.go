package properties

// Corners sets a different radius per corner (CSS border-radius with
// four values), in dp. When any corner is non-zero it takes precedence
// over BorderRadius.
type Corners struct {
	TopLeft, TopRight, BottomRight, BottomLeft int
}

// UniformCorners returns corners with the same radius.
func UniformCorners(r int) Corners {
	return Corners{TopLeft: r, TopRight: r, BottomRight: r, BottomLeft: r}
}

// TopCorners rounds only the two top corners.
func TopCorners(r int) Corners {
	return Corners{TopLeft: r, TopRight: r}
}

// BottomCorners rounds only the two bottom corners.
func BottomCorners(r int) Corners {
	return Corners{BottomRight: r, BottomLeft: r}
}

// IsZero reports whether no corner is set.
func (c Corners) IsZero() bool {
	return c == Corners{}
}

// ResolveCorners returns the per-corner radii for an element: c when
// set, otherwise the uniform radius r on every corner.
func ResolveCorners(r BorderRadius, c Corners) Corners {
	if !c.IsZero() {
		return c
	}
	return UniformCorners(int(r))
}

// clampCorners keeps every radius within half the box size so the
// rounded rectangle stays well formed.
func clampCorners(c Corners, w, h int) Corners {
	lim := min(w, h) / 2
	cl := func(v int) int {
		if v < 0 {
			return 0
		}
		return min(v, lim)
	}
	return Corners{cl(c.TopLeft), cl(c.TopRight), cl(c.BottomRight), cl(c.BottomLeft)}
}
