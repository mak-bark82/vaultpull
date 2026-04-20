package envrollback

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRestore_NilSnapshot(t *testing.T) {
	_, err := Restore(".env", nil)
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestRestore_EmptyTargetFile(t *testing.T) {
	snap := &Snapshot{Data: map[string]string{"K": "V"}}
	_, err := Restore("", snap)
	if err == nil {
		t.Fatal("expected error for empty target file")
	}
}

func TestRestore_WritesKeys(t *testing.T) {
	target := filepath.Join(t.TempDir(), ".env")
	snap := &Snapshot{
		File:      target,
		Timestamp: time.Now(),
		Data:      map[string]string{"HOST": "db", "PORT": "5432"},
	}
	res, err := Restore(target, snap)
	if err != nil {
		t.Fatalf("Restore failed: %v", err)
	}
	if res.Written != 2 {
		t.Errorf("expected 2 written, got %d", res.Written)
	}
	raw, _ := os.ReadFile(target)
	content := string(raw)
	if !strings.Contains(content, "HOST=db") {
		t.Errorf("expected HOST=db in output, got: %s", content)
	}
	if !strings.Contains(content, "PORT=5432") {
		t.Errorf("expected PORT=5432 in output, got: %s", content)
	}
}

func TestRestore_OverwritesExistingFile(t *testing.T) {
	target := filepath.Join(t.TempDir(), ".env")
	_ = os.WriteFile(target, []byte("OLD=value\n"), 0600)
	snap := &Snapshot{
		Data: map[string]string{"NEW": "fresh"},
	}
	_, err := Restore(target, snap)
	if err != nil {
		t.Fatalf("Restore failed: %v", err)
	}
	raw, _ := os.ReadFile(target)
	if strings.Contains(string(raw), "OLD") {
		t.Errorf("old content should have been removed")
	}
	if !strings.Contains(string(raw), "NEW=fresh") {
		t.Errorf("expected NEW=fresh in restored file")
	}
}

func TestRestore_EmptySnapshot(t *testing.T) {
	target := filepath.Join(t.TempDir(), ".env")
	snap := &Snapshot{Data: map[string]string{}}
	res, err := Restore(target, snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Written != 0 {
		t.Errorf("expected 0 written, got %d", res.Written)
	}
}
