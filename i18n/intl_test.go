package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIntl(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "locales")
	os.MkdirAll(dir, 0o755)
	os.WriteFile(filepath.Join(dir, "en.json"), []byte(`{
		"hello": "Hello, {name}!",
		"files": {"one": "{count} file in {folder}", "other": "{count} files in {folder}"}
	}`), 0o644)
	os.WriteFile(filepath.Join(dir, "de.json"), []byte(`{"hello": "Hallo, {name}!"}`), 0o644)
	if err := LoadDir(dir); err != nil {
		t.Fatal(err)
	}
	defer SetLocale(Locale())
	SetLocale("en")
	type user struct {
		Name string
		ID   int `json:"id"`
	}
	for _, tc := range []struct {
		got, want string
	}{
		{Intl("hello", map[string]any{"name": "Ana"}), "Hello, Ana!"},
		{Intl("hello", map[string]string{"name": "Ana"}), "Hello, Ana!"},
		{Intl("hello", user{Name: "Ana"}), "Hello, Ana!"},
		{Intl("hello", &user{Name: "Ana"}), "Hello, Ana!"},
		{Intl("hello", "name", "Ana"), "Hello, Ana!"},
		{Intl("files", map[string]any{"count": 1, "folder": "Download"}), "1 file in Download"},
		{Intl("files", map[string]any{"count": int64(3), "folder": "Download"}), "3 files in Download"},
		{Intl("missing {x}", "x", 1), "missing 1"},
	} {
		if tc.got != tc.want {
			t.Errorf("got %q, want %q", tc.got, tc.want)
		}
	}
	SetLocale("de-AT")
	if got := Intl("hello", "name", "Ana"); got != "Hallo, Ana!" {
		t.Errorf("de-AT: %q", got)
	}
	if got := Intl("files", "count", 2, "folder", "x"); got != "2 files in x" {
		t.Errorf("fallback to en: %q", got)
	}
	if err := LoadDir(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Error("a missing folder loaded")
	}
}
