// Package bridge helps giode's platform packages (push, store, review)
// call into Java and Objective-C: asynchronous requests answered by a
// callback with an id, and, on Android, static Java calls.
package bridge

import (
	"context"
	"sync"
)

// Requests pairs asynchronous native requests with their answers.
type Requests[T any] struct {
	mu   sync.Mutex
	next int
	m    map[int]chan T
}

// New starts a request: the id to pass to native code, and the channel
// Done answers on.
func (r *Requests[T]) New() (int, chan T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.m == nil {
		r.m = map[int]chan T{}
	}
	r.next++
	ch := make(chan T, 1)
	r.m[r.next] = ch
	return r.next, ch
}

// Done answers request id; unknown or finished ids are ignored.
func (r *Requests[T]) Done(id int, v T) {
	r.mu.Lock()
	ch := r.m[id]
	delete(r.m, id)
	r.mu.Unlock()
	if ch != nil {
		ch <- v
	}
}

// Cancel forgets request id.
func (r *Requests[T]) Cancel(id int) {
	r.mu.Lock()
	delete(r.m, id)
	r.mu.Unlock()
}

// Wait waits for request id's answer, or for ctx.
func (r *Requests[T]) Wait(ctx context.Context, id int, ch chan T) (T, error) {
	select {
	case v := <-ch:
		return v, nil
	case <-ctx.Done():
		r.Cancel(id)
		var zero T
		return zero, ctx.Err()
	}
}

// AfterNotifier, when set, runs after giode creates its system
// notifier (App.Notify), which on macOS takes over the notification
// center's delegate.
var AfterNotifier func()
