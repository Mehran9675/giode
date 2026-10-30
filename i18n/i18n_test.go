package i18n

import (
	"os"
	"runtime"
	"testing"

	"gioui.org/io/system"
)

func TestTranslate(t *testing.T) {
	if err := Load(os.DirFS("testdata"), "."); err != nil {
		t.Fatal(err)
	}
	changes := 0
	OnChange(func() { changes++ })
	SetLocale("en_US.UTF-8")
	if Locale() != "en-us" || RTL() {
		t.Errorf("locale %q rtl %v", Locale(), RTL())
	}
	if got := T("greeting", "name", "Ana"); got != "Hello, Ana!" {
		t.Errorf("T = %q", got)
	}
	if got := N("unread", 1); got != "1 message" {
		t.Errorf("N(1) = %q", got)
	}
	if got := N("unread", 5); got != "5 messages" {
		t.Errorf("N(5) = %q", got)
	}
	SetLocale("fa-IR")
	if !RTL() || Direction() != system.RTL {
		t.Error("Persian isn't right to left")
	}
	if got := T("greeting", map[string]any{"name": "سارا"}); got != "سلام، سارا!" {
		t.Errorf("fa T = %q", got)
	}
	if got := T("only_en"); got != "English only" {
		t.Errorf("fallback = %q", got)
	}
	if got := T("missing.key"); got != "missing.key" {
		t.Errorf("missing = %q", got)
	}
	SetLocale("ru")
	for n, want := range map[int]string{1: "1 сообщение", 3: "3 сообщения", 5: "5 сообщений", 11: "11 сообщений", 21: "21 сообщение", 22: "22 сообщения"} {
		if got := N("unread", n); got != want {
			t.Errorf("ru N(%d) = %q, want %q", n, got, want)
		}
	}
	if changes != 3 {
		t.Errorf("%d change notifications, want 3", changes)
	}
}

func TestPluralRules(t *testing.T) {
	cases := []struct {
		lang string
		n    int
		want string
	}{
		{"en", 0, "other"}, {"en", 1, "one"}, {"fr", 0, "one"}, {"fr", 2, "other"},
		{"ar", 0, "zero"}, {"ar", 2, "two"}, {"ar", 5, "few"}, {"ar", 15, "many"}, {"ar", 100, "other"},
		{"pl", 22, "few"}, {"pl", 25, "many"}, {"ja", 1, "other"}, {"he", 2, "two"},
		{"lv", 0, "zero"}, {"lv", 1, "one"}, {"lv", 2, "other"}, {"lv", 11, "zero"}, {"lv", 21, "one"}, {"lv", 30, "zero"}, {"lv", 111, "zero"},
		{"sl", 1, "one"}, {"sl", 2, "two"}, {"sl", 3, "few"}, {"sl", 4, "few"}, {"sl", 5, "other"}, {"sl", 101, "one"}, {"sl", 102, "two"}, {"sl", 104, "few"}, {"sl", 11, "other"},
		{"is", 21, "one"}, {"is", 11, "other"}, {"ga", 4, "few"}, {"ga", 8, "many"}, {"cy", 6, "many"}, {"mt", 13, "many"},
	}
	for _, c := range cases {
		if got := pluralForm(c.lang, c.n); got != c.want {
			t.Errorf("%s %d: %s, want %s", c.lang, c.n, got, c.want)
		}
	}
}

func TestSystemLocale(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "de_DE.UTF-8")
	got := SystemLocale()
	switch runtime.GOOS {
	case "darwin", "windows", "ios", "android", "js":
		// These read the system's language setting, not $LANG.
		if got == "" || got != normalize(got) {
			t.Errorf("SystemLocale = %q, not a normalized tag", got)
		}
	default:
		if got != "de-de" && got != "en" {
			t.Errorf("SystemLocale = %q", got)
		}
	}
}
