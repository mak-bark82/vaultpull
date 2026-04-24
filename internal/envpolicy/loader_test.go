package envpolicy

import (
	"os"
	"path/filepath"
	"testing"
)

func writeRulesFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "policy.yaml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeRulesFile: %v", err)
	}
	return p
}

func TestLoadRules_Valid(t *testing.T) {
	p := writeRulesFile(t, `
rules:
  - name: no-debug
    pattern: "^DEBUG"
    target: key
    action: deny
  - name: warn-plain
    pattern: "^hello"
    target: value
    action: warn
`)
	rules, err := LoadRules(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
}

func TestLoadRules_EmptyPath(t *testing.T) {
	rules, err := LoadRules("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected 0 rules, got %d", len(rules))
	}
}

func TestLoadRules_MissingFile(t *testing.T) {
	_, err := LoadRules("/nonexistent/policy.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadRules_MissingName(t *testing.T) {
	p := writeRulesFile(t, `
rules:
  - pattern: "^DEBUG"
    target: key
    action: deny
`)
	_, err := LoadRules(p)
	if err == nil {
		t.Fatal("expected error for missing rule name")
	}
}

func TestLoadRules_UnknownAction(t *testing.T) {
	p := writeRulesFile(t, `
rules:
  - name: bad-action
    pattern: "^DEBUG"
    target: key
    action: block
`)
	_, err := LoadRules(p)
	if err == nil {
		t.Fatal("expected error for unknown action")
	}
}

func TestLoadRules_DefaultsTargetToKey(t *testing.T) {
	p := writeRulesFile(t, `
rules:
  - name: no-debug
    pattern: "^DEBUG"
    action: warn
`)
	rules, err := LoadRules(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules[0].Target != "key" {
		t.Errorf("expected default target 'key', got %q", rules[0].Target)
	}
}
