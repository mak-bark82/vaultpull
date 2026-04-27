package envencrypt

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "encrypt.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return p
}

func TestLoadConfig_EmptyPath(t *testing.T) {
	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Patterns) != 0 {
		t.Errorf("expected empty patterns, got %v", cfg.Patterns)
	}
}

func TestLoadConfig_Valid(t *testing.T) {
	path := writeConfigFile(t, "patterns:\n  - SECRET\n  - TOKEN\n")
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if len(cfg.Patterns) != 2 {
		t.Errorf("expected 2 patterns, got %d", len(cfg.Patterns))
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/encrypt.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfig_NoPatterns(t *testing.T) {
	path := writeConfigFile(t, "patterns: []\n")
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error when patterns is empty")
	}
}

func TestLoadConfig_EmptyPatternEntry(t *testing.T) {
	path := writeConfigFile(t, "patterns:\n  - SECRET\n  - \"\"\n")
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for empty pattern entry")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	path := writeConfigFile(t, ": bad: yaml: [\n")
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
