// Package state provides a minimal global-state primitive: an
// observable value that can be read and written from anywhere - not
// just from within a Layout call, but a goroutine outside the render
// loop too (a network response, a timer) - and requests a redraw when
// it changes.
//
// Reads are tracked: a cached part of the UI (giode.Memo) that read a
// State while it was built is rebuilt when that State changes, and
// only then.
package state

import (
	"reflect"
	"sync"
	"sync/atomic"
)

// invalidate requests a redraw of the app's window. giode's own
// App.New sets this once; state doesn't import the root giode
// package itself (giode imports state, not the other way around), so
// this is how the two stay connected.
var invalidate atomic.Pointer[func()]

// SetInvalidator registers the function State calls after a value
// changes. The app does this once, in App.New; application code using
// State doesn't need to call it.
func SetInvalidator(fn func()) {
	if fn == nil {
		invalidate.Store(nil)
		return
	}
	invalidate.Store(&fn)
}

// RequestRedraw asks the app for a new frame. It is safe to call from
// any goroutine; components that change on their own outside a frame
// (a toast shown from a background goroutine) use it.
func RequestRedraw() {
	if fn := invalidate.Load(); fn != nil {
		(*fn)()
	}
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
	order    []int
	nextID   int
	version  atomic.Uint64
	name     string
}

// New returns a State holding initial.
func New[T any](initial T) *State[T] {
	return &State[T]{value: initial, watchers: map[int]func(T){}}
}

// Named returns a State holding initial, registered under name so
// tooling - the built-in MCP server's state tools - can read and set
// it. Names must be unique; registering a name again replaces the
// earlier State.
func Named[T any](name string, initial T) *State[T] {
	s := New(initial)
	s.name = name
	register(name, s)
	return s
}

// Get returns the current value.
func (s *State[T]) Get() T {
	s.mu.RLock()
	v := s.value
	s.mu.RUnlock()
	recordRead(s, s.version.Load())
	return v
}

// Set replaces the value, notifies every watcher (in registration
// order, on the calling goroutine), and requests a redraw.
func (s *State[T]) Set(v T) {
	s.mu.Lock()
	s.value = v
	s.version.Add(1)
	s.mu.Unlock()
	s.notify(v)
}

// SetIfChanged is Set, skipping the update - no watchers, no redraw,
// no cache invalidation - when v equals the current value. Values that
// aren't comparable with == are compared with reflect.DeepEqual.
func (s *State[T]) SetIfChanged(v T) bool {
	s.mu.Lock()
	if equal(s.value, v) {
		s.mu.Unlock()
		return false
	}
	s.value = v
	s.version.Add(1)
	s.mu.Unlock()
	s.notify(v)
	return true
}

// Update atomically replaces the value with fn(current), notifies
// every watcher and requests a redraw - safe for a read-modify-write
// under concurrent access (e.g. incrementing a counter) where a
// separate Get then Set could race with another goroutine's write in
// between.
func (s *State[T]) Update(fn func(current T) T) {
	s.mu.Lock()
	s.value = fn(s.value)
	s.version.Add(1)
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
	s.order = append(s.order, id)
	s.mu.Unlock()
	return func() {
		s.mu.Lock()
		delete(s.watchers, id)
		for i, o := range s.order {
			if o == id {
				s.order = append(s.order[:i], s.order[i+1:]...)
				break
			}
		}
		s.mu.Unlock()
	}
}

// Version returns a counter incremented by every change.
func (s *State[T]) Version() uint64 {
	return s.version.Load()
}

func (s *State[T]) currentVersion() uint64 {
	return s.version.Load()
}

// Name returns the name the State was registered under with Named, or
// "".
func (s *State[T]) Name() string {
	return s.name
}

// notify calls every watcher with v, in registration order, then
// requests a redraw. Watchers are copied out under the lock first so a
// watcher that calls Set, Watch or its own unsubscribe function
// doesn't deadlock or race the map being ranged over.
func (s *State[T]) notify(v T) {
	s.mu.RLock()
	watchers := make([]func(T), 0, len(s.order))
	for _, id := range s.order {
		watchers = append(watchers, s.watchers[id])
	}
	s.mu.RUnlock()

	for _, w := range watchers {
		w(v)
	}
	if fn := invalidate.Load(); fn != nil {
		(*fn)()
	}
}

func equal[T any](a, b T) bool {
	if t := reflect.TypeOf(a); t != nil && t.Comparable() {
		return any(a) == any(b)
	}
	return reflect.DeepEqual(a, b)
}
