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
