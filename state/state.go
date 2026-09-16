// Package state provides a minimal global-state primitive: an
// observable value that can be read and written from anywhere - not
// just from within a Layout call, but a goroutine outside the render
// loop too (a network response, a timer) - and requests a redraw when
// it changes.
package state

import "sync"

// invalidate requests a redraw of the app's window. giode's own
// App.New sets this once; state doesn't import the root giode
// package itself (giode imports state, not the other way around), so
// this is how the two stay connected.
var invalidate func()

// SetInvalidator registers the function State calls after a value
// changes. The app does this once, in App.New; application code using
// State doesn't need to call it.
func SetInvalidator(fn func()) {
	invalidate = fn
}

// State is an observable value of type T, safe for concurrent use:
// read and write it from anywhere, including a goroutine outside the
// render loop, not just from within a Layout call. Create it once,
// typically as a package-level var or a field alongside a stateful
// component's own state, and share the pointer wherever the value is
// read or written.
type State[T any] struct {
	mu       sync.RWMutex
	value    T
	watchers map[int]func(T)
	nextID   int
}

// New returns a State holding initial.
func New[T any](initial T) *State[T] {
	return &State[T]{value: initial, watchers: map[int]func(T){}}
}

// Get returns the current value.
func (s *State[T]) Get() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

// Set replaces the value, notifies every watcher (in registration
// order, on the calling goroutine), and requests a redraw.
func (s *State[T]) Set(v T) {
	s.mu.Lock()
	s.value = v
	s.mu.Unlock()
	s.notify(v)
}

// Update atomically replaces the value with fn(current), notifies
// every watcher and requests a redraw - safe for a read-modify-write
// under concurrent access (e.g. incrementing a counter) where a
// separate Get then Set could race with another goroutine's write in
// between.
func (s *State[T]) Update(fn func(current T) T) {
	s.mu.Lock()
	s.value = fn(s.value)
	v := s.value
	s.mu.Unlock()
	s.notify(v)
}

// Watch registers fn to be called every time Set/Update changes the
// value. It returns an unsubscribe function.
func (s *State[T]) Watch(fn func(T)) (unsubscribe func()) {
	s.mu.Lock()
	id := s.nextID
	s.nextID++
	s.watchers[id] = fn
	s.mu.Unlock()
	return func() {
		s.mu.Lock()
		delete(s.watchers, id)
		s.mu.Unlock()
	}
}

// notify calls every watcher with v, then requests a redraw. Watchers
// are copied out under the lock first so a watcher that calls Set,
// Watch or its own unsubscribe function doesn't deadlock or race the
// map being ranged over.
func (s *State[T]) notify(v T) {
	s.mu.RLock()
	watchers := make([]func(T), 0, len(s.watchers))
	for _, w := range s.watchers {
		watchers = append(watchers, w)
	}
	s.mu.RUnlock()

	for _, w := range watchers {
		w(v)
	}
	if invalidate != nil {
		invalidate()
	}
}
