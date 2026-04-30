package envclassify_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nahK994/vaultpull/internal/envclassify"
)

func writeRulesFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "classify.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write rules file: %v", err)
	}
	return p
}

func TestLoadRules_EmptyPath(t *testing.T) {
	rules, err := envclassify.LoadRules("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Fatalf("expected empty rules")
	}
}

func TestLoadRules_MissingFile(t *testing.T) {
	_, err := envclassify.LoadRules("/no/such/file.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadRules_Valid(t *testing.T) {
	content := `rules:
  - pattern: "(?i)secret"
    category: secret
  - pattern: "^DB_"
    category: database
`
	p := writeRulesFile(t, content)
	rules, err := envclassify.LoadRules(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
}

func TestLoadRules_MissingPattern(t *testing.T) {
	content := `rules:
  - pattern: ""
    category: secret
`
	p := writeRulesFile(t, content)
	_, err := envclassify.LoadRules(p)
	if err == nil {
		t.Fatal("expected error for missing pattern")
	}
}

func TestLoadRules_MissingCategory(t *testing.T) {
	content := `rules:
  - pattern: "SECRET"
    category: ""
`
	p := writeRulesFile(t, content)
	_, err := envclassify.LoadRules(p)
	if err == nil {
		t.Fatal("expected error for missing category")
	}
}
