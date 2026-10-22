package state

import (
	"sync"
	"sync/atomic"
)

// Dependency tracking.
//
// Cached UI (giode's Memo) records which States its subtree read while
// it was built, then skips rebuilding until one of them changes. Every
// State carries a version counter for this; reads made while a Tracker
// is active are recorded with the version they saw.

// versioned is implemented by every State.
type versioned interface {
	currentVersion() uint64
}

// Tracker records the States read while it is active. Trackers nest:
// a read is recorded by every active tracker, since a nested cache's
// dependencies are also its parents'.
type Tracker struct {
	deps map[versioned]uint64
}

var (
	trackMu  sync.Mutex
	trackers []*Tracker
	// tracking is the number of active trackers, so Get can skip the
	// lock entirely when nothing is being tracked - the common case
	// outside cached layout.
	tracking atomic.Int32
)

// Track starts recording State reads into a new Tracker. Call Stop
// when done. Reads from any goroutine are recorded while a tracker is
// active; a spurious dependency only costs an extra rebuild.
func Track() *Tracker {
	t := &Tracker{}
	trackMu.Lock()
	trackers = append(trackers, t)
	trackMu.Unlock()
	tracking.Add(1)
	return t
}

// Stop ends recording. Trackers must be stopped in reverse order of
// creation.
func (t *Tracker) Stop() {
	trackMu.Lock()
	for i := len(trackers) - 1; i >= 0; i-- {
		if trackers[i] == t {
			trackers = append(trackers[:i], trackers[i+1:]...)
			break
		}
	}
	trackMu.Unlock()
	tracking.Add(-1)
}

// Changed reports whether any State read while t was active has
// changed since. A nil Tracker never changes.
func (t *Tracker) Changed() bool {
	if t == nil {
		return false
	}
	for s, v := range t.deps {
		if s.currentVersion() != v {
			return true
		}
	}
	return false
}

// Len returns the number of distinct States recorded.
func (t *Tracker) Len() int {
	if t == nil {
		return 0
	}
	return len(t.deps)
}

// recordRead registers a read of s at version v with every active
// tracker.
func recordRead(s versioned, v uint64) {
	if tracking.Load() == 0 {
		return
	}
	trackMu.Lock()
	for _, t := range trackers {
		if t.deps == nil {
			t.deps = map[versioned]uint64{}
		}
		if _, ok := t.deps[s]; !ok {
			t.deps[s] = v
		}
	}
	trackMu.Unlock()
}
