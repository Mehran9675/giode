package properties

// Inset is the space around an element: Top, Right, Bottom, Left.
// Start and End add to the side where lines start and end, like CSS's
// padding-inline-start and -end: the left and right, swapped in a
// right-to-left locale. Inset{Start: 16} indents from the start in
// either direction.
type Inset struct {
	Top, Right, Bottom, Left int
	Start, End               int
}

// Physical resolves Start and End into Left and Right for a layout
// reading right to left or not.
func (i Inset) Physical(rtl bool) Inset {
	if i.Start == 0 && i.End == 0 {
		return i
	}
	if rtl {
		i.Right += i.Start
		i.Left += i.End
	} else {
		i.Left += i.Start
		i.Right += i.End
	}
	i.Start, i.End = 0, 0
	return i
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
