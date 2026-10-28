package properties

// ZIndex controls the stacking order of siblings (CSS z-index): higher
// values paint above and receive input before lower ones. It applies to
// flow and absolutely positioned children alike; at equal z-index a
// positioned child paints above a flow child, and otherwise children
// paint in document order. It doesn't change where children are laid
// out.
type ZIndex int
