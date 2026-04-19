package envwriter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWrite_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	secrets := map[string]string{"APP_KEY": "abc123", "DB_URL": "postgres://localhost/db"}
	if err := Write(path, secrets, Options{Overwrite: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "APP_KEY=abc123") {
		t.Errorf("expected APP_KEY=abc123 in output, got:\n%s", content)
	}
}

func TestWrite_NoOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	_ = os.WriteFile(path, []byte("APP_KEY=original\n"), 0600)

	secrets := map[string]string{"APP_KEY": "new_value", "NEW_KEY": "hello"}
	if err := Write(path, secrets, Options{Overwrite: false}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "APP_KEY=original") {
		t.Errorf("expected APP_KEY to remain original, got:\n%s", content)
	}
	if !strings.Contains(content, "NEW_KEY=hello") {
		t.Errorf("expected NEW_KEY=hello to be added, got:\n%s", content)
	}
}

func TestWrite_Overwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	_ = os.WriteFile(path, []byte("APP_KEY=original\n"), 0600)

	secrets := map[string]string{"APP_KEY": "updated"}
	if err := Write(path, secrets, Options{Overwrite: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), "APP_KEY=updated") {
		t.Errorf("expected APP_KEY=updated, got:\n%s", string(data))
	}
}
