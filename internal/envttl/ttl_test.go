package envttl_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/envttl"
)

func tempStore(t *testing.T) *envttl.Store {
	t.Helper()
	dir := t.TempDir()
	s, err := envttl.New(filepath.Join(dir, "ttl.json"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestNew_EmptyPath(t *testing.T) {
	_, err := envttl.New("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestSet_And_IsExpired_NotYet(t *testing.T) {
	s := tempStore(t)
	if err := s.Set("DB_PASS", 10*time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if s.IsExpired("DB_PASS") {
		t.Error("expected key to not be expired")
	}
}

func TestIsExpired_UnknownKey(t *testing.T) {
	s := tempStore(t)
	if s.IsExpired("UNKNOWN") {
		t.Error("unknown key should not be considered expired")
	}
}

func TestSet_EmptyKey_ReturnsError(t *testing.T) {
	s := tempStore(t)
	if err := s.Set("", time.Minute); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestExpired_ReturnsExpiredKeys(t *testing.T) {
	dir := t.TempDir()
	s, _ := envttl.New(filepath.Join(dir, "ttl.json"))

	// Manually set a past TTL by writing directly then reloading via a new store
	_ = s.Set("FRESH", 10*time.Minute)
	_ = s.Set("STALE", -1*time.Second) // already elapsed

	expired := s.Expired()
	if len(expired) != 1 || expired[0] != "STALE" {
		t.Errorf("expected [STALE], got %v", expired)
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	s := tempStore(t)
	_ = s.Set("TOKEN", time.Hour)
	if err := s.Remove("TOKEN"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if s.IsExpired("TOKEN") {
		t.Error("removed key should not be expired")
	}
}

func TestPersistence_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ttl.json")

	s1, _ := envttl.New(path)
	_ = s1.Set("API_KEY", 5*time.Minute)

	s2, err := envttl.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if s2.IsExpired("API_KEY") {
		t.Error("reloaded key should not be expired")
	}
}

func TestNew_LoadsExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ttl.json")
	// write invalid JSON to confirm error surface
	_ = os.WriteFile(path, []byte("not-json"), 0o600)
	_, err := envttl.New(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
