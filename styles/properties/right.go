package properties

// Right is the offset in dp from the right edge of the parent's content
// box for absolutely positioned elements. It anchors the element's
// right edge when Left is unset; with both set, a box without a Width
// stretches between them. 0 means unset; use Flush for an offset of
// zero.
type Right int
