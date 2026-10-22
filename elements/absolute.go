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

// AbsolutePosition returns the window position of the element being
// laid out, as far as it is known at that moment. Only meaningful
// during the current frame's layout.
//
// Containers that place a child only after laying it out (a Row or
// Stack with several children, a Scroll, Tabs) don't know the child's
// offset yet, so inside such a child the value is the container's
// origin. Anything recorded from it through RegisterPortal,
// RegisterPopup, RecordSemantic or TrackPosition is moved to the right
// place once the container has placed the child (see ShiftPositions);
// a value used directly in the same Layout call is not.
func AbsolutePosition() image.Point {
	return currentOffset
}

// PushOffset adds off to the tracked position until the returned
// function runs; pair it with the op.Offset around the same child. An
// element that knows a child's offset before laying it out calls this
// so AbsolutePosition is exact inside the child; one that only knows
// it afterwards uses MarkPositions and ShiftPositions instead.
func PushOffset(off image.Point) (pop func()) {
	currentOffset = currentOffset.Add(off)
	return func() { currentOffset = currentOffset.Sub(off) }
}
