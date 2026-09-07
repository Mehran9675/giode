package properties

// Inset is the space around an element: Top, Right, Bottom, Left.
type Inset struct {
	Top, Right, Bottom, Left int
}

// UniformInset returns an inset with equal sides.
func UniformInset(v int) Inset {
	return Inset{Top: v, Right: v, Bottom: v, Left: v}
}

// SymmetricInset returns an inset symmetric in the horizontal and
// vertical axes.
func SymmetricInset(h, v int) Inset {
	return Inset{Top: v, Right: h, Bottom: v, Left: h}
}
