package elements

import (
	"math"
	"sync"
	"sync/atomic"

	"gioui.org/unit"
)

// Threading.
//
// The view function, every Layout, scroll content closures and event
// handlers (OnClick, OnChange, ...) all run on one goroutine: the UI
// goroutine. Data changed only there needs no locks. Code running
// elsewhere - a network callback, a timer, an MCP tool - hands work to
// the UI goroutine with App.Post (or App.RunOnUI to wait for it), or
// uses State, which is safe from any goroutine.

var (
	afterMu     sync.Mutex
	afterLayout []func()
)

// AfterLayout runs fn once the current frame's layout is finished, on
// the UI goroutine, then requests another frame. A handler that must
// change the model the rest of this frame is still reading - removing
// the list item being laid out, say - defers the change with it.
// Called outside a frame, fn runs after the next one.
func AfterLayout(fn func()) {
	afterMu.Lock()
	afterLayout = append(afterLayout, fn)
	afterMu.Unlock()
}

// RunAfterLayout runs and clears the AfterLayout queue, reporting
// whether anything ran. The app calls it after laying out each frame;
// component authors don't need to.
func RunAfterLayout() bool {
	afterMu.Lock()
	fns := afterLayout
	afterLayout = nil
	afterMu.Unlock()
	for _, fn := range fns {
		fn()
	}
	if len(fns) > 0 {
		InvalidateCaches()
	}
	return len(fns) > 0
}

// displayScale is the last frame's pixels per dp, stored as float32
// bits so it can be read from any goroutine.
var displayScale atomic.Uint32

func init() {
	displayScale.Store(math.Float32bits(1))
}

// SetMetric records the frame's display metric. The app calls it every
// frame; component authors don't need to.
func SetMetric(m unit.Metric) {
	s := m.PxPerDp
	if s <= 0 {
		s = 1
	}
	displayScale.Store(math.Float32bits(s))
}

// DisplayScale returns the display's pixels per dp as of the last
// frame (1 on a standard display, 2 on a typical HiDPI one). It is
// safe to call from anywhere, including outside Layout, to convert
// sizes for code that doesn't have a layout context.
func DisplayScale() float32 {
	return math.Float32frombits(displayScale.Load())
}

// DpToPx converts dp to pixels with the last frame's display scale.
func DpToPx(v float32) int {
	return int(math.Round(float64(v * DisplayScale())))
}

// PxToDp converts pixels to dp with the last frame's display scale.
func PxToDp(px int) float32 {
	return float32(px) / DisplayScale()
}
