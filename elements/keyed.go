package elements

import (
	"sync"
	"weak"
)

// Keyed state.
//
// Stateful components (a Pressable, a text field, a Transition) must
// survive across frames, which is easy for a fixed UI - create them
// once - but not for a list whose items come and go at runtime. Keyed
// gives each item its own persistent state, looked up by a stable key
// (an ID, not an index), and forgets state nobody asked for in a while,
// so removed items don't leak.

// DefaultRetainFrames is how many frames an entry may go unrequested
// before it's dropped. Frames only happen when something changes, so
// this is not a fixed time.
const DefaultRetainFrames = 60

type keyedKey struct{ scope, key any }

type keyedEntry struct {
	value any
	frame uint64
}

var (
	keyedMu     sync.Mutex
	keyedStore  = map[keyedKey]*keyedEntry{}
	keyedRetain = uint64(DefaultRetainFrames)
	// keyedCaches sweeps every live KeyedCache; each func reports
	// false once its cache was garbage collected.
	keyedCaches  []func(frame uint64) bool
	keyedCacheMu sync.Mutex
)

// Keyed returns the state stored under (scope, key), creating it with
// init the first time it's asked for. Use it inside the view function
// to give each item of a dynamic list its own stateful component:
//
//	for _, f := range files {
//		row := giode.Keyed("files", f.ID, func() *giode.PressableElement {
//			return giode.Pressable(...)
//		})
//		rows = append(rows, row)
//	}
//
// scope separates unrelated lists that might share keys; any
// comparable values work for both. If a key was stored with a
// different type, init replaces it. Entries not asked for in
// DefaultRetainFrames frames are dropped (see SetKeyedRetention). Call
// it on the UI goroutine.
func Keyed[T any](scope, key any, init func() T) T {
	k := keyedKey{scope, key}
	keyedMu.Lock()
	e, ok := keyedStore[k]
	if ok {
		if v, ok := e.value.(T); ok {
			e.frame = frameNumber
			keyedMu.Unlock()
			retainKeyed(e.touch)
			return v
		}
	}
	keyedMu.Unlock()
	v := init()
	e = &keyedEntry{value: v, frame: frameNumber}
	keyedMu.Lock()
	keyedStore[k] = e
	keyedMu.Unlock()
	retainKeyed(e.touch)
	return v
}

// touch keeps the entry alive this frame.
func (e *keyedEntry) touch() {
	keyedMu.Lock()
	e.frame = frameNumber
	keyedMu.Unlock()
}

// retainStack collects, for caches recording their subtree, the keyed
// lookups to repeat on frames they replay instead.
var retainStack []*[]func()

func pushRetainScope(s *[]func()) (pop func()) {
	retainStack = append(retainStack, s)
	n := len(retainStack)
	return func() { retainStack = retainStack[:n-1] }
}

// retainKeyed registers touch with the innermost recording cache.
func retainKeyed(touch func()) {
	if n := len(retainStack); n > 0 {
		s := retainStack[n-1]
		*s = append(*s, touch)
	}
}

// ForgetKeyed drops the state stored under (scope, key) right away.
func ForgetKeyed(scope, key any) {
	keyedMu.Lock()
	delete(keyedStore, keyedKey{scope, key})
	keyedMu.Unlock()
}

// SetKeyedRetention sets how many frames Keyed entries survive without
// being asked for (minimum 1).
func SetKeyedRetention(frames int) {
	keyedMu.Lock()
	keyedRetain = uint64(max(1, frames))
	keyedMu.Unlock()
}

// sweepKeyed drops expired entries; EndFrame calls it.
func sweepKeyed() {
	keyedMu.Lock()
	for k, e := range keyedStore {
		if frameNumber-e.frame > keyedRetain {
			delete(keyedStore, k)
		}
	}
	keyedMu.Unlock()
	keyedCacheMu.Lock()
	live := keyedCaches[:0]
	for _, sweep := range keyedCaches {
		if sweep(frameNumber) {
			live = append(live, sweep)
		}
	}
	clear(keyedCaches[len(live):])
	keyedCaches = live
	keyedCacheMu.Unlock()
}

// KeyedCache is a typed, self-contained version of Keyed: a map of
// per-item state that drops entries not asked for in a while. Hold one
// in a stateful component:
//
//	type fileList struct {
//		rows *giode.KeyedCache[string, *giode.PressableElement]
//	}
//	...
//	row := l.rows.Get(file.ID, func() *giode.PressableElement { ... })
//
// Expired entries are dropped at the end of each frame. It is safe for
// concurrent use, though state is normally only touched on the UI
// goroutine.
type KeyedCache[K comparable, V any] struct {
	mu      sync.Mutex
	entries map[K]*keyedCacheEntry[V]
	retain  uint64
	onEvict func(K, V)
}

type keyedCacheEntry[V any] struct {
	value V
	frame uint64
}

// NewKeyedCache returns an empty cache that keeps entries for
// DefaultRetainFrames unrequested frames.
func NewKeyedCache[K comparable, V any]() *KeyedCache[K, V] {
	c := &KeyedCache[K, V]{entries: map[K]*keyedCacheEntry[V]{}, retain: DefaultRetainFrames}
	wp := weak.Make(c)
	keyedCacheMu.Lock()
	keyedCaches = append(keyedCaches, func(frame uint64) bool {
		if c := wp.Value(); c != nil {
			c.sweep(frame)
			return true
		}
		return false
	})
	keyedCacheMu.Unlock()
	return c
}

// Retain sets how many frames an entry survives without being asked
// for (minimum 1). It returns c for chaining.
func (c *KeyedCache[K, V]) Retain(frames int) *KeyedCache[K, V] {
	c.mu.Lock()
	c.retain = uint64(max(1, frames))
	c.mu.Unlock()
	return c
}

// OnEvict sets a function called with each dropped entry, to release
// resources it holds. It returns c for chaining.
func (c *KeyedCache[K, V]) OnEvict(fn func(K, V)) *KeyedCache[K, V] {
	c.mu.Lock()
	c.onEvict = fn
	c.mu.Unlock()
	return c
}

// Get returns the entry for key, creating it with init the first time.
func (c *KeyedCache[K, V]) Get(key K, init func() V) V {
	c.mu.Lock()
	if e, ok := c.entries[key]; ok {
		e.frame = frameNumber
		c.mu.Unlock()
		retainKeyed(func() { c.touch(e) })
		return e.value
	}
	c.mu.Unlock()
	v := init()
	e := &keyedCacheEntry[V]{value: v, frame: frameNumber}
	c.mu.Lock()
	c.entries[key] = e
	c.mu.Unlock()
	retainKeyed(func() { c.touch(e) })
	return v
}

func (c *KeyedCache[K, V]) touch(e *keyedCacheEntry[V]) {
	c.mu.Lock()
	e.frame = frameNumber
	c.mu.Unlock()
}

// Peek returns the entry for key without creating it or keeping it
// alive.
func (c *KeyedCache[K, V]) Peek(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		var zero V
		return zero, false
	}
	return e.value, true
}

// Delete drops the entry for key right away.
func (c *KeyedCache[K, V]) Delete(key K) {
	c.mu.Lock()
	e, ok := c.entries[key]
	delete(c.entries, key)
	fn := c.onEvict
	c.mu.Unlock()
	if ok && fn != nil {
		fn(key, e.value)
	}
}

// Len returns the number of entries.
func (c *KeyedCache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}

// Range calls fn for each entry until it returns false.
func (c *KeyedCache[K, V]) Range(fn func(K, V) bool) {
	c.mu.Lock()
	snapshot := make(map[K]V, len(c.entries))
	for k, e := range c.entries {
		snapshot[k] = e.value
	}
	c.mu.Unlock()
	for k, v := range snapshot {
		if !fn(k, v) {
			return
		}
	}
}

func (c *KeyedCache[K, V]) sweep(frame uint64) {
	type evicted struct {
		k K
		v V
	}
	var gone []evicted
	c.mu.Lock()
	for k, e := range c.entries {
		if frame-e.frame > c.retain {
			gone = append(gone, evicted{k, e.value})
			delete(c.entries, k)
		}
	}
	fn := c.onEvict
	c.mu.Unlock()
	if fn != nil {
		for _, g := range gone {
			fn(g.k, g.v)
		}
	}
}
