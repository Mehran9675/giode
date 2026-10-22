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
	orig := invalidate.Load()
	defer invalidate.Store(orig)

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

func TestState_WatchersRunInRegistrationOrder(t *testing.T) {
	s := New(0)
	var order []int
	for i := 0; i < 20; i++ {
		i := i
		s.Watch(func(int) { order = append(order, i) })
	}
	s.Set(1)
	for i, v := range order {
		if v != i {
			t.Fatalf("watcher order = %v, want ascending", order)
		}
	}
}

func TestState_SetIfChanged(t *testing.T) {
	s := New([]int{1, 2})
	v := s.Version()
	if s.SetIfChanged([]int{1, 2}) {
		t.Error("SetIfChanged with an equal slice reported a change")
	}
	if s.Version() != v {
		t.Error("version advanced without a change")
	}
	if !s.SetIfChanged([]int{3}) {
		t.Error("SetIfChanged with a different value reported no change")
	}
}

func TestTracker_RecordsReads(t *testing.T) {
	a, b := New(1), New(2)
	tr := Track()
	_ = a.Get()
	tr.Stop()
	_ = b.Get() // after Stop: not tracked
	if tr.Len() != 1 {
		t.Fatalf("tracked %d states, want 1", tr.Len())
	}
	if tr.Changed() {
		t.Error("Changed before any write")
	}
	b.Set(3)
	if tr.Changed() {
		t.Error("write to an untracked state reported as a change")
	}
	a.Set(5)
	if !tr.Changed() {
		t.Error("write to a tracked state not reported")
	}
}

func TestNamed_JSONRoundTrip(t *testing.T) {
	type prefs struct {
		Volume int    `json:"volume"`
		Theme  string `json:"theme"`
	}
	s := Named("test.prefs", prefs{Volume: 3, Theme: "dark"})
	defer Unregister("test.prefs")
	e, ok := Lookup("test.prefs")
	if !ok {
		t.Fatal("named state not registered")
	}
	data, err := e.JSON()
	if err != nil || string(data) != `{"volume":3,"theme":"dark"}` {
		t.Fatalf("JSON = %s, %v", data, err)
	}
	if err := e.SetJSON([]byte(`{"volume":7,"theme":"light"}`)); err != nil {
		t.Fatal(err)
	}
	if got := s.Get(); got.Volume != 7 || got.Theme != "light" {
		t.Errorf("after SetJSON: %+v", got)
	}
	if err := e.SetJSON([]byte(`"nope"`)); err == nil {
		t.Error("SetJSON with the wrong type: want error")
	}
}
