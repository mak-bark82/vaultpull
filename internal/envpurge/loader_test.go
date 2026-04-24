package envpurge

import (
	"os"
	"path/filepath"
	"testing"
)

func writeRulesFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "purge_rules.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadRules_Valid(t *testing.T) {
	p := writeRulesFile(t, "rules:\n  - key: DB_PASSWORD\n    reason: sensitive\n  - key: API_KEY\n    reason: rotated\n")
	rules, err := LoadRules(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Key != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD, got %s", rules[0].Key)
	}
}

func TestLoadRules_EmptyPath(t *testing.T) {
	rules, err := LoadRules("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected empty rules")
	}
}

func TestLoadRules_MissingFile(t *testing.T) {
	_, err := LoadRules("/nonexistent/purge.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadRules_MissingKey(t *testing.T) {
	p := writeRulesFile(t, "rules:\n  - reason: no key here\n")
	_, err := LoadRules(p)
	if err == nil {
		t.Error("expected error for rule missing 'key'")
	}
}
