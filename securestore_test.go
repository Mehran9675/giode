package giode

import (
	"errors"
	"os"
	"testing"
)

// TestSecureStore runs against the platform's credential store. On
// Linux it needs a Secret Service (GIODE_TEST_KEYRING=1, under a
// session with an unlocked GNOME Keyring); elsewhere it skips when the
// store isn't reachable.
func TestSecureStore(t *testing.T) {
	if os.Getenv("GIODE_TEST_KEYRING") == "" {
		t.Skip("set GIODE_TEST_KEYRING=1 to run against the real credential store")
	}
	s := &SecureStorage{Service: "org.giode.test"}
	if err := s.Set("api key", []byte("sk-secret")); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := s.Set("api key", []byte("sk-replaced")); err != nil {
		t.Fatalf("Set again: %v", err)
	}
	v, err := s.Get("api key")
	if err != nil || string(v) != "sk-replaced" {
		t.Fatalf("Get = %q, %v; want the replaced value", v, err)
	}
	if err := s.Delete("api key"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := s.Get("api key"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get after Delete: %v, want ErrNotFound", err)
	}
	if err := s.Delete("api key"); err != nil {
		t.Errorf("Delete of a missing name: %v", err)
	}
}
