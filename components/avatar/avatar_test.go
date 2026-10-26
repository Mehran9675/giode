package avatar

import "testing"

func TestInitialsAndColor(t *testing.T) {
	for in, want := range map[string]string{
		"Ada Lovelace":           "AL",
		"grace":                  "G",
		"Jean-Luc Picard":        "JP",
		"  John Ronald Tolkien ": "JT",
		"":                       "",
	} {
		if got := Initials(in); got != want {
			t.Errorf("Initials(%q) = %q, want %q", in, got, want)
		}
	}
	if Color("Ada") != Color("ada") {
		t.Error("color depends on case")
	}
}
