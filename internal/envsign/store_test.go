package envSign_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	envSign "github.com/yourusername/vaultpull/internal/envSign"
)

func TestSaveAndLoadRecord_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sig.json")

	now := time.Now().UTC().Truncate(time.Second)
	rec := envSign.Record{
		File:      ".env",
		Signature: "abc123",
		SignedAt:  now,
	}

	if err := envSign.SaveRecord(path, rec); err != nil {
		t.Fatalf("SaveRecord failed: %v", err)
	}

	loaded, err := envSign.LoadRecord(path)
	if err != nil {
		t.Fatalf("LoadRecord failed: %v", err)
	}

	if loaded.File != rec.File {
		t.Errorf("File mismatch: got %q, want %q", loaded.File, rec.File)
	}
	if loaded.Signature != rec.Signature {
		t.Errorf("Signature mismatch: got %q, want %q", loaded.Signature, rec.Signature)
	}
	if !loaded.SignedAt.Equal(rec.SignedAt) {
		t.Errorf("SignedAt mismatch: got %v, want %v", loaded.SignedAt, rec.SignedAt)
	}
}

func TestLoadRecord_EmptyPath(t *testing.T) {
	_, err := envSign.LoadRecord("")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestLoadRecord_MissingFile(t *testing.T) {
	_, err := envSign.LoadRecord("/nonexistent/path/sig.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadRecord_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0600)

	_, err := envSign.LoadRecord(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestSaveRecord_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sig.json")

	rec := envSign.Record{File: ".env", Signature: "xyz", SignedAt: time.Now()}
	if err := envSign.SaveRecord(path, rec); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to be created")
	}
}
