package envreader

import (
	"os"
	"testing"
)

func writeFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envreader_*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestRead_ValidFile(t *testing.T) {
	path := writeFile(t, "FOO=bar\nBAZ=qux\n")
	m, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if m["FOO"] != "bar" || m["BAZ"] != "qux" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestRead_SkipsComments(t *testing.T) {
	path := writeFile(t, "# comment\nFOO=bar\n")
	m, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 1 || m["FOO"] != "bar" {
		t.Errorf("unexpected map: %v", m)
	}
}

func TestRead_NonExistentFile(t *testing.T) {
	m, err := Read("/tmp/does_not_exist_vaultpull.env")
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestRead_SkipsMalformedLines(t *testing.T) {
	path := writeFile(t, "NOEQUALS\nFOO=bar\n")
	m, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 1 {
		t.Errorf("expected 1 entry, got %v", m)
	}
}
