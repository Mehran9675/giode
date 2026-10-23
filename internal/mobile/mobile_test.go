package mobile

import "testing"

func TestLifecycleHandlers(t *testing.T) {
	var bg, fg int
	SetLifecycleHandlers(func() { bg++ }, func() { fg++ })
	notify(&onBackground)
	notify(&onForeground)
	notify(&onForeground)
	if bg != 1 || fg != 2 {
		t.Errorf("bg=%d fg=%d, want 1 2", bg, fg)
	}
	SetLifecycleHandlers(nil, nil)
	notify(&onBackground)
	if bg != 1 {
		t.Error("cleared handler still ran")
	}
}
