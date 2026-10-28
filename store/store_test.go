package store

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	ps, err := parseProducts(`[{"id":"pro","title":"Pro","price":"€4,99","micros":4990000,"currency":"EUR"},{"id":"club","subscription":true,"period":"P1M"}]`)
	if err != nil || len(ps) != 2 {
		t.Fatal(ps, err)
	}
	if ps[0].PriceMicros != 4990000 || ps[0].Currency != "EUR" || ps[0].Price != "€4,99" || ps[0].Subscription {
		t.Errorf("%+v", ps[0])
	}
	if !ps[1].Subscription || ps[1].Period != "P1M" {
		t.Errorf("%+v", ps[1])
	}
	p, err := parsePurchase(`{"id":"GPA.1","product":"pro","token":"tok","time":1700000000000,"pending":true,"receipt":"{}","signature":"sig"}`)
	if err != nil || p.ProductID != "pro" || p.Token != "tok" || !p.Pending || p.Signature != "sig" || !p.Time.Equal(time.UnixMilli(1700000000000)) {
		t.Errorf("%+v %v", p, err)
	}
	if nativeError("canceled") != ErrCanceled || nativeError("owned") != ErrOwned || nativeError("") != nil || nativeError("x").Error() != "store: x" {
		t.Error("native errors")
	}
}

func TestUpdatesQueueUntilOnUpdate(t *testing.T) {
	defer OnUpdate(nil)
	update(Purchase{ID: "early"})
	got := make(chan Purchase, 2)
	OnUpdate(func(p Purchase) { got <- p })
	update(Purchase{ID: "late"})
	seen := map[string]bool{}
	for i := 0; i < 2; i++ {
		select {
		case p := <-got:
			seen[p.ID] = true
		case <-time.After(time.Second):
			t.Fatal("missing update")
		}
	}
	if !seen["early"] || !seen["late"] {
		t.Errorf("%v", seen)
	}
}
