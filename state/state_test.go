package state

import (
	"sync"
	"testing"
)

func TestState_GetSet(t *testing.T) {
	s := New(1)
	if got := s.Get(); got != 1 {
		t.Fatalf("Get() = %d, want 1", got)
	}
	s.Set(2)
	if got := s.Get(); got != 2 {
		t.Fatalf("Get() after Set(2) = %d, want 2", got)
	}
}

func TestState_Update(t *testing.T) {
	s := New(10)
	s.Update(func(v int) int { return v + 5 })
	if got := s.Get(); got != 15 {
		t.Fatalf("Get() after Update(+5) = %d, want 15", got)
	}
}

func TestState_WatchFiresOnSetAndUpdate(t *testing.T) {
	s := New("a")
	var seen []string
	s.Watch(func(v string) { seen = append(seen, v) })

	s.Set("b")
	s.Update(func(v string) string { return v + "!" })

	want := []string{"b", "b!"}
	if len(seen) != len(want) {
		t.Fatalf("watcher saw %v, want %v", seen, want)
	}
	for i := range want {
		if seen[i] != want[i] {
			t.Errorf("watcher call %d = %q, want %q", i, seen[i], want[i])
		}
	}
}

func TestState_UnsubscribeStopsFutureNotifications(t *testing.T) {
	s := New(0)
	calls := 0
	unsubscribe := s.Watch(func(int) { calls++ })

	s.Set(1)
	if calls != 1 {
		t.Fatalf("calls = %d after one Set, want 1", calls)
	}

	unsubscribe()
	s.Set(2)
	if calls != 1 {
		t.Fatalf("calls = %d after unsubscribing and a further Set, want 1 (unchanged)", calls)
	}
}

func TestState_MultipleWatchersAllFire(t *testing.T) {
	s := New(0)
	var a, b int
	s.Watch(func(v int) { a = v })
	s.Watch(func(v int) { b = v })

	s.Set(42)
	if a != 42 || b != 42 {
		t.Errorf("a=%d b=%d after Set(42), want both 42", a, b)
	}
}

func TestState_SetCallsInvalidator(t *testing.T) {
	orig := invalidate
	defer func() { invalidate = orig }()

	calls := 0
	SetInvalidator(func() { calls++ })

	s := New(0)
	s.Set(1)
	s.Update(func(v int) int { return v + 1 })

	if calls != 2 {
		t.Errorf("invalidator called %d times, want 2 (one per Set/Update)", calls)
	}
}

// TestState_ConcurrentAccess exercises Get/Set/Update/Watch from many
// goroutines at once under the race detector (go test -race), the
// main thing a "safe for concurrent use" primitive needs to prove.
func TestState_ConcurrentAccess(t *testing.T) {
	s := New(0)
	unsubscribe := s.Watch(func(int) {})
	defer unsubscribe()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			s.Set(i)
		}()
		go func() {
			defer wg.Done()
			s.Update(func(v int) int { return v + 1 })
		}()
		go func() {
			defer wg.Done()
			_ = s.Get()
		}()
	}
	wg.Wait()
}
