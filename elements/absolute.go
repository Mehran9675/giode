package elements

import "image"

// currentOffset accumulates the absolute position of the current
// point in a Layout traversal. Layout is single-threaded and
// synchronous (no concurrent goroutines walk the tree), so a single
// package-level accumulator, pushed before laying out a child and
// popped after, correctly tracks "where am I right now" for anything
// that needs its own absolute screen position - the building block
// behind kit.RegisterPortal (and therefore Tooltip, nested Menu
// submenus and AutoComplete's panel).
var currentOffset image.Point

// AbsolutePosition returns the accumulated offset at the current
// point in the Layout traversal: where the calling element's origin
// currently sits relative to the window. Only meaningful while being
// laid out as part of the current frame.
//
// Known limitation: it's exact for a Box's own content offset (margin
// and padding) and for a left/top-anchored absolutely positioned
// child, which is the common nesting case. It's an approximation - the
// shared container's origin, not the individual child's own offset -
// for one of several children distributed by a flex container (a Row
// or Stack with more than one child): a flex container measures every
// child before it knows where any of them land, so by the time a
// child's final position is known, that child (and anything nested in
// it) has already finished its own Layout call. Closing this requires
// instrumenting the flex distribution itself, which real Gio flex
// containers don't expose partway through; tracked as a follow-up.
func AbsolutePosition() image.Point {
	return currentOffset
}

// PushOffset accumulates off into the tracked absolute position;
// pair with the op.Offset(off).Push()/Pop() around the same child. A
// custom element that positions its own children with op.Offset
// should call this alongside it so AbsolutePosition (and therefore
// kit.RegisterPortal) works correctly for content placed inside it.
func PushOffset(off image.Point) (pop func()) {
	currentOffset = currentOffset.Add(off)
	return func() { currentOffset = currentOffset.Sub(off) }
}
