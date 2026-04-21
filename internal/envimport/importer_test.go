package envimport_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/envimport"
)

func writeSrc(t *testing.T, content string) string {
	t.Helper()
	tmp := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(tmp, []byte(content), 0o600); err != nil {
		t.Fatalf("writeSrc: %v", err)
	}
	return tmp
}

func TestImport_AddsNewKeys(t *testing.T) {
	src := writeSrc(t, "FOO=bar\nBAZ=qux\n")
	dst := map[string]string{}
	result, err := envimport.Import(src, dst, envimport.Options{Format: envimport.FormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "bar" || result["BAZ"] != "qux" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestImport_NoOverwrite_PreservesExisting(t *testing.T) {
	src := writeSrc(t, "FOO=new\n")
	dst := map[string]string{"FOO": "original"}
	result, err := envimport.Import(src, dst, envimport.Options{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "original" {
		t.Errorf("expected 'original', got %q", result["FOO"])
	}
}

func TestImport_Overwrite_ReplacesExisting(t *testing.T) {
	src := writeSrc(t, "FOO=new\n")
	dst := map[string]string{"FOO": "original"}
	result, err := envimport.Import(src, dst, envimport.Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "new" {
		t.Errorf("expected 'new', got %q", result["FOO"])
	}
}

func TestImport_SkipsCommentsAndBlanks(t *testing.T) {
	src := writeSrc(t, "# comment\n\nKEY=val\n")
	result, err := envimport.Import(src, map[string]string{}, envimport.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result["KEY"] != "val" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestImport_MissingFile_ReturnsError(t *testing.T) {
	_, err := envimport.Import("/nonexistent/.env", map[string]string{}, envimport.Options{})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestImport_UnsupportedFormat_ReturnsError(t *testing.T) {
	src := writeSrc(t, "KEY=val\n")
	_, err := envimport.Import(src, map[string]string{}, envimport.Options{Format: "toml"})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestImport_QuotedValues(t *testing.T) {
	src := writeSrc(t, `KEY="hello world"`+"\n")
	result, err := envimport.Import(src, map[string]string{}, envimport.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", result["KEY"])
	}
}
