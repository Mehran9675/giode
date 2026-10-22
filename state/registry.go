package state

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

// Entry is a named State seen through a type-erased view, for tooling
// that inspects or edits application state (the MCP server).
type Entry interface {
	// Name is the registered name.
	Name() string
	// JSON returns the current value encoded as JSON.
	JSON() ([]byte, error)
	// SetJSON decodes data into a value of the State's type and sets
	// it.
	SetJSON(data []byte) error
}

var (
	registryMu sync.RWMutex
	registry   = map[string]Entry{}
)

func register[T any](name string, s *State[T]) {
	registryMu.Lock()
	registry[name] = entry[T]{s}
	registryMu.Unlock()
}

// Lookup returns the named State registered with Named.
func Lookup(name string) (Entry, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	e, ok := registry[name]
	return e, ok
}

// Entries returns every named State, sorted by name.
func Entries() []Entry {
	registryMu.RLock()
	out := make([]Entry, 0, len(registry))
	for _, e := range registry {
		out = append(out, e)
	}
	registryMu.RUnlock()
	sort.Slice(out, func(i, j int) bool { return out[i].Name() < out[j].Name() })
	return out
}

// Unregister removes a named State from the registry.
func Unregister(name string) {
	registryMu.Lock()
	delete(registry, name)
	registryMu.Unlock()
}

type entry[T any] struct{ s *State[T] }

func (e entry[T]) Name() string { return e.s.name }

func (e entry[T]) JSON() ([]byte, error) {
	return json.Marshal(e.s.Get())
}

func (e entry[T]) SetJSON(data []byte) error {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("state %q: %w", e.s.name, err)
	}
	e.s.Set(v)
	return nil
}
