package giode

import (
	"testing"

	"gioui.org/io/key"
	"gioui.org/io/pointer"
)

func TestParseKeyName(t *testing.T) {
	for in, want := range map[string]key.Name{
		"enter": key.NameReturn, "Tab": key.NameTab, "page_down": key.NamePageDown, "backspace": key.NameDeleteBackward,
		"a": "A", "Z": "Z", "5": "5", "/": "/", "f1": key.NameF1, "F12": key.NameF12, "left": key.NameLeftArrow,
	} {
		got, err := parseKeyName(in)
		if err != nil || got != want {
			t.Errorf("parseKeyName(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, in := range []string{"f13", "f01", "nope", ""} {
		if _, err := parseKeyName(in); err == nil {
			t.Errorf("parseKeyName(%q) succeeded", in)
		}
	}
}

func TestParseModifiers(t *testing.T) {
	m, err := parseModifiers([]string{"ctrl", "Shift", "shortcut"})
	if err != nil || m != key.ModCtrl|key.ModShift|key.ModShortcut {
		t.Errorf("parseModifiers = %v, %v", m, err)
	}
	if _, err := parseModifiers([]string{"hyper"}); err == nil {
		t.Error("unknown modifier accepted")
	}
	if b, err := parseButton("right"); err != nil || b != pointer.ButtonSecondary {
		t.Errorf("parseButton(right) = %v, %v", b, err)
	}
}
