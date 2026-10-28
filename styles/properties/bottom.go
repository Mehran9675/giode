package properties

// Bottom is the offset in dp from the bottom edge of the parent's
// content box for absolutely positioned elements. It anchors the
// element's bottom edge when Top is unset; with both set, a box without
// a Height stretches between them. 0 means unset; use Flush for an
// offset of zero.
type Bottom int
