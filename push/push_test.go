package push

import (
	"testing"
	"time"
)

func TestParseMessage(t *testing.T) {
	m, err := parseMessage(`{"id":"m1","title":"Hi","body":"There","data":{"chat":"42","n":7,"o":{"a":true}}}`)
	if err != nil {
		t.Fatal(err)
	}
	if m.ID != "m1" || m.Title != "Hi" || m.Body != "There" {
		t.Errorf("%+v", m)
	}
	if m.Data["chat"] != "42" || m.Data["n"] != "7" || m.Data["o"] != `{"a":true}` {
		t.Errorf("data %v", m.Data)
	}
	if _, err := parseMessage("nope"); err == nil {
		t.Error("bad JSON parsed")
	}
}

// TestQueueUntilListen checks events from before Listen (the
// notification that launched the app) reach the handlers.
func TestQueueUntilListen(t *testing.T) {
	defer func() { started, queued, handlers = false, nil, Handlers{} }()
	received(`{"title":"Opened"}`, true)
	deliver(event{kind: evToken, token: "tok"})

	opened := make(chan Message, 1)
	tokens := make(chan string, 1)
	msgs := make(chan Message, 1)
	Listen(nil, Handlers{
		Open:    func(m Message) { opened <- m },
		Token:   func(s string) { tokens <- s },
		Message: func(m Message) { msgs <- m },
	})
	if m := <-opened; m.Title != "Opened" {
		t.Errorf("opened %+v", m)
	}
	if s := <-tokens; s != "tok" {
		t.Errorf("token %q", s)
	}
	received(`{"body":"live"}`, false)
	select {
	case m := <-msgs:
		if m.Body != "live" {
			t.Errorf("message %+v", m)
		}
	case <-time.After(time.Second):
		t.Error("no message after Listen")
	}
}
