package envpin_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envpin"
)

func writePinFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writePinFile: %v", err)
	}
	return path
}

func TestLoad_EmptyPath(t *testing.T) {
	p, err := envpin.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Entries()) != 0 {
		t.Errorf("expected no entries, got %d", len(p.Entries()))
	}
}

func TestLoad_NonExistent(t *testing.T) {
	p, err := envpin.Load("/nonexistent/pins.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.Entries()) != 0 {
		t.Errorf("expected no entries")
	}
}

func TestPin_And_IsPinned(t *testing.T) {
	p, _ := envpin.Load("")
	p.Pin("DB_PASSWORD", "secret123", "locked for prod")

	if !p.IsPinned("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be pinned")
	}
	if p.IsPinned("OTHER_KEY") {
		t.Error("OTHER_KEY should not be pinned")
	}
}

func TestUnpin_RemovesEntry(t *testing.T) {
	p, _ := envpin.Load("")
	p.Pin("API_KEY", "abc", "")
	removed := p.Unpin("API_KEY")
	if !removed {
		t.Error("expected Unpin to return true")
	}
	if p.IsPinned("API_KEY") {
		t.Error("API_KEY should no longer be pinned")
	}
	if p.Unpin("API_KEY") {
		t.Error("second Unpin should return false")
	}
}

func TestApply_OverridesWithPinnedValue(t *testing.T) {
	p, _ := envpin.Load("")
	p.Pin("DB_PASS", "pinned_value", "")

	incoming := map[string]string{
		"DB_PASS": "vault_value",
		"OTHER":   "unchanged",
	}
	out := p.Apply(incoming)

	if out["DB_PASS"] != "pinned_value" {
		t.Errorf("expected pinned_value, got %q", out["DB_PASS"])
	}
	if out["OTHER"] != "unchanged" {
		t.Errorf("expected unchanged, got %q", out["OTHER"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	p, _ := envpin.Load("")
	p.Pin("X", "pinned", "")

	incoming := map[string]string{"X": "original"}
	p.Apply(incoming)

	if incoming["X"] != "original" {
		t.Error("Apply must not mutate the input map")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")

	p, _ := envpin.Load("")
	p.Pin("SECRET", "val", "important")
	if err := p.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	p2, err := envpin.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !p2.IsPinned("SECRET") {
		t.Error("expected SECRET to be pinned after reload")
	}
	entries := p2.Entries()
	if entries[0].Comment != "important" {
		t.Errorf("expected comment 'important', got %q", entries[0].Comment)
	}
}
