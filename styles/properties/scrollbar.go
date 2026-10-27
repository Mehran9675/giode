package properties

// ScrollBar customizes the scrollbar a Scroll container draws when
// its content overflows. Zero values mean "unset" and fall back to
// the built-in defaults.
type ScrollBar struct {
	// Width is the thickness of the bar in pixels (default 8).
	Width int
	// TrackColor is the fill behind the thumb (default a subtle,
	// translucent gray).
	TrackColor Color
	// ThumbColor is the draggable indicator (default a light gray).
	ThumbColor Color
	// Radius rounds the track and thumb corners in pixels (default
	// half the resolved Width, giving a pill shape).
	Radius int
	// MinThumbLength is the smallest the thumb may shrink to on very
	// long content, so it stays visible and draggable (default 20).
	MinThumbLength int

	// Hidden hides the bar. The content still scrolls with the wheel,
	// touch, the keyboard and ScrollBy.
	Hidden bool
	// HiddenMaxWidth hides the bar while the window is at most this
	// many dp wide - a phone layout - and shows it in wider windows.
	HiddenMaxWidth int
	// HiddenOn hides the bar while the innermost Responsive matched a
	// breakpoint with one of these names, separated by commas (for
	// example "compact" or "compact,medium").
	HiddenOn string
}
