package envrewrite

import (
	"os"
	"path/filepath"
	"testing"
)

func writeRulesFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "rules.yaml")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write rules file: %v", err)
	}
	return p
}

func TestLoadRules_Valid(t *testing.T) {
	p := writeRulesFile(t, `
rules:
  - key: APP_HOST
    find: localhost
    replace: prod.example.com
    target: value
  - find: OLD_
    replace: NEW_
    target: key
`)
	rs, err := LoadRules(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rs.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rs.Rules))
	}
	if rs.Rules[0].Find != "localhost" {
		t.Errorf("unexpected find value: %s", rs.Rules[0].Find)
	}
}

func TestLoadRules_EmptyPath(t *testing.T) {
	rs, err := LoadRules("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rs.Rules) != 0 {
		t.Error("expected empty rules for empty path")
	}
}

func TestLoadRules_MissingFile(t *testing.T) {
	_, err := LoadRules("/nonexistent/rules.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadRules_MissingFind(t *testing.T) {
	p := writeRulesFile(t, `
rules:
  - key: APP_HOST
    replace: newvalue
    target: value
`)
	_, err := LoadRules(p)
	if err == nil {
		t.Error("expected error for missing find field")
	}
}

func TestLoadRules_InvalidTarget(t *testing.T) {
	p := writeRulesFile(t, `
rules:
  - find: old
    replace: new
    target: invalid
`)
	_, err := LoadRules(p)
	if err == nil {
		t.Error("expected error for invalid target")
	}
}
