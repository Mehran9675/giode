package properties

// Visibility fades the element and its children, in the range 0
// (fully hidden) to 1 (fully visible). A zero value means unset and
// resolves to fully visible, mirroring Opacity; the two combine
// multiplicatively.
type Visibility float32
