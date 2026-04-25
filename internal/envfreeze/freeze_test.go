package envfreeze_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/envfreeze"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "freeze.json")
}

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_EmptyPath(t *testing.T) {
	_, err := envfreeze.New("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNew_NonExistent(t *testing.T) {
	f, err := envfreeze.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Keys()) != 0 {
		t.Fatal("expected empty keys for new store")
	}
}

func TestFreeze_AddsKeys(t *testing.T) {
	f, _ := envfreeze.New(tempPath(t))
	ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := f.Freeze([]string{"DB_PASS", "API_KEY"}, "initial", fixedClock(ts)); err != nil {
		t.Fatalf("Freeze: %v", err)
	}
	if !f.IsFrozen("DB_PASS") {
		t.Error("expected DB_PASS to be frozen")
	}
	if !f.IsFrozen("API_KEY") {
		t.Error("expected API_KEY to be frozen")
	}
}

func TestFreeze_Persists(t *testing.T) {
	path := tempPath(t)
	f1, _ := envfreeze.New(path)
	_ = f1.Freeze([]string{"SECRET"}, "", nil)

	f2, err := envfreeze.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if !f2.IsFrozen("SECRET") {
		t.Error("expected SECRET to survive reload")
	}
}

func TestIsFrozen_Unknown(t *testing.T) {
	f, _ := envfreeze.New(tempPath(t))
	if f.IsFrozen("MISSING") {
		t.Error("expected MISSING to not be frozen")
	}
}

func TestUnfreeze_RemovesKey(t *testing.T) {
	f, _ := envfreeze.New(tempPath(t))
	_ = f.Freeze([]string{"A", "B", "C"}, "", nil)
	if err := f.Unfreeze("B"); err != nil {
		t.Fatalf("Unfreeze: %v", err)
	}
	if f.IsFrozen("B") {
		t.Error("expected B to be unfrozen")
	}
	if !f.IsFrozen("A") || !f.IsFrozen("C") {
		t.Error("expected A and C to remain frozen")
	}
}

func TestFreeze_SkipsEmptyKeys(t *testing.T) {
	f, _ := envfreeze.New(tempPath(t))
	_ = f.Freeze([]string{"", "VALID", ""}, "", nil)
	keys := f.Keys()
	if len(keys) != 1 || keys[0] != "VALID" {
		t.Errorf("expected [VALID], got %v", keys)
	}
}

func TestNew_InvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := envfreeze.New(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
