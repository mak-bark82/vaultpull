package envcompare_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envcompare"
)

func writeEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	return p
}

func TestLoadFile_Valid(t *testing.T) {
	dir := t.TempDir()
	p := writeEnvFile(t, dir, ".env", "FOO=bar\nBAZ=\"quoted\"\n# comment\n\nBAD_LINE\n")
	env, err := envcompare.LoadFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["FOO"] != "bar" {
		t.Errorf("FOO: got %q", env["FOO"])
	}
	if env["BAZ"] != "quoted" {
		t.Errorf("BAZ: got %q", env["BAZ"])
	}
	if _, ok := env["BAD_LINE"]; ok {
		t.Error("malformed line should be skipped")
	}
}

func TestLoadFile_NonExistent(t *testing.T) {
	_, err := envcompare.LoadFile("/no/such/file.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestCompareFiles_Differences(t *testing.T) {
	dir := t.TempDir()
	left := writeEnvFile(t, dir, "left.env", "A=1\nB=old\n")
	right := writeEnvFile(t, dir, "right.env", "B=new\nC=3\n")
	r, err := envcompare.CompareFiles(left, right)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.OnlyInLeft["A"]; !ok {
		t.Error("expected A only-left")
	}
	if _, ok := r.OnlyInRight["C"]; !ok {
		t.Error("expected C only-right")
	}
	if _, ok := r.Different["B"]; !ok {
		t.Error("expected B in Different")
	}
}

func TestCompareFiles_MissingLeft(t *testing.T) {
	dir := t.TempDir()
	right := writeEnvFile(t, dir, "right.env", "K=v\n")
	_, err := envcompare.CompareFiles("/no/left.env", right)
	if err == nil {
		t.Error("expected error for missing left file")
	}
}
