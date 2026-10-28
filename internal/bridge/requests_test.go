package bridge

import (
	"context"
	"testing"
	"time"
)

func TestRequests(t *testing.T) {
	var r Requests[string]
	id, ch := r.New()
	go r.Done(id, "ok")
	if v, err := r.Wait(context.Background(), id, ch); v != "ok" || err != nil {
		t.Fatalf("%q %v", v, err)
	}
	r.Done(id, "again") // finished: ignored, doesn't block

	id, ch = r.New()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if _, err := r.Wait(ctx, id, ch); err != context.DeadlineExceeded {
		t.Fatalf("timeout: %v", err)
	}
	r.Done(id, "late") // canceled: ignored
}
