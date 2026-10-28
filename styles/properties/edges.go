package properties

// Flush is an explicit zero for Top, Right, Bottom and Left, whose zero
// value means "not set": Bottom: properties.Flush pins an absolutely
// positioned box to the bottom edge, where Bottom: 0 would leave it at
// the top.
const Flush = -1 << 31

// edge converts an offset to pixels: unset and Flush are both 0.
func edge(v int) int {
	if v == Flush {
		return 0
	}
	return v
}

// IsSet reports whether the offset was set (Flush counts).
func (t Top) IsSet() bool { return t != 0 }

// Px returns the offset in pixels, 0 for Flush.
func (t Top) Px() int { return edge(int(t)) }

// IsSet reports whether the offset was set (Flush counts).
func (r Right) IsSet() bool { return r != 0 }

// Px returns the offset in pixels, 0 for Flush.
func (r Right) Px() int { return edge(int(r)) }

// IsSet reports whether the offset was set (Flush counts).
func (b Bottom) IsSet() bool { return b != 0 }

// Px returns the offset in pixels, 0 for Flush.
func (b Bottom) Px() int { return edge(int(b)) }

// IsSet reports whether the offset was set (Flush counts).
func (l Left) IsSet() bool { return l != 0 }

// Px returns the offset in pixels, 0 for Flush.
func (l Left) Px() int { return edge(int(l)) }
