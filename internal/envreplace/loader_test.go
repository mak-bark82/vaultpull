package envreplace_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envreplace"
)

func writeRulesFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "rules.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write rules file: %v", err)
	}
	return p
}

func TestLoadRules_Valid(t *testing.T) {
	path := writeRulesFile(t, `
rules:
  - pattern: "localhost"
    with: "db.prod.internal"
  - key: API_URL
    pattern: "http"
    with: "https"
`)
	rules, err := envreplace.LoadRules(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[1].Key != "API_URL" {
		t.Errorf("expected key API_URL, got %q", rules[1].Key)
	}
}

func TestLoadRules_EmptyPath(t *testing.T) {
	rules, err := envreplace.LoadRules("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected empty rules, got %d", len(rules))
	}
}

func TestLoadRules_MissingFile(t *testing.T) {
	_, err := envreplace.LoadRules("/nonexistent/rules.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadRules_MissingPattern(t *testing.T) {
	path := writeRulesFile(t, `
rules:
  - key: FOO
    with: bar
`)
	_, err := envreplace.LoadRules(path)
	if err == nil {
		t.Fatal("expected error for missing pattern")
	}
}

func TestLoadRules_InvalidYAML(t *testing.T) {
	path := writeRulesFile(t, `{invalid yaml: [`)
	_, err := envreplace.LoadRules(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
