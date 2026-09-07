package properties

// ZIndex controls the stacking order of siblings (CSS z-index).
// Children with higher values are laid out later: they paint above
// and receive input before siblings with lower values. Unlike CSS,
// z-index also affects the layout order in flex containers.
type ZIndex int
