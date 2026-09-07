package properties

// Width is the width of an element in pixels.
//
//   - 0 (default): hug the content, up to the available width
//   - -1: expand to fill the available width
//   - >0: fixed width
type Width int

// Constrain applies the width to a pair of minimum and maximum
// constraints along an axis.
func (w Width) Constrain(min, max int) (int, int) {
	switch {
	case w > 0:
		return int(w), int(w)
	case w == -1 && max < unbounded:
		return max, max
	default:
		return min, max
	}
}

// unbounded mirrors the sentinel Gio uses for infinite constraints.
const unbounded = 1e6
