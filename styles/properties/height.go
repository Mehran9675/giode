package properties

// Height is the height of an element in pixels.
//
//   - 0 (default): hug the content, up to the available height
//   - -1: expand to fill the available height
//   - >0: fixed height
type Height int

// Constrain applies the height to a pair of minimum and maximum
// constraints along an axis.
func (h Height) Constrain(min, max int) (int, int) {
	switch {
	case h > 0:
		return int(h), int(h)
	case h == -1 && max < unbounded:
		return max, max
	default:
		return min, max
	}
}
