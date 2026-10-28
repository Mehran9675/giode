package giode

import (
	"net/url"
	"runtime"
	"testing"
)

func TestOpenURL_SecondInstanceAndStartup(t *testing.T) {
	var got []string
	a := New(Config{OnOpenURL: func(u *url.URL) { got = append(got, u.String()) }})
	a.handleLaunch([]string{"myapp://chat/42?x=1"})
	a.runUICalls()
	if len(got) != 1 || got[0] != "myapp://chat/42?x=1" {
		t.Errorf("second instance link: %v", got)
	}
	got = nil
	a.startupURL([]string{"myapp://start"})
	a.runUICalls()
	if runtime.GOOS == "linux" && (len(got) != 1 || got[0] != "myapp://start") {
		t.Errorf("startup link: %v", got)
	}
	got = nil
	a.handleLaunch([]string{"notes.txt"})
	a.runUICalls()
	if len(got) != 0 {
		t.Errorf("a plain argument opened %v", got)
	}
}
