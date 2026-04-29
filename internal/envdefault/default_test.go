package envdefault_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/envdefault"
)

func writeRulesFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "defaults.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write rules file: %v", err)
	}
	return p
}

func TestApply_FillsMissingKey(t *testing.T) {
	rules := []envdefault.Rule{{Key: "LOG_LEVEL", Default: "info"}}
	out, results, err := envdefault.Apply(map[string]string{}, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["LOG_LEVEL"] != "info" {
		t.Errorf("expected 'info', got %q", out["LOG_LEVEL"])
	}
	if len(results) != 1 || results[0].Reason != "missing" {
		t.Errorf("expected one missing result, got %+v", results)
	}
}

func TestApply_SkipsExistingKey(t *testing.T) {
	rules := []envdefault.Rule{{Key: "LOG_LEVEL", Default: "info"}}
	out, results, err := envdefault.Apply(map[string]string{"LOG_LEVEL": "debug"}, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["LOG_LEVEL"] != "debug" {
		t.Errorf("expected 'debug', got %q", out["LOG_LEVEL"])
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %+v", results)
	}
}

func TestApply_OnEmpty_FillsEmptyValue(t *testing.T) {
	rules := []envdefault.Rule{{Key: "TIMEOUT", Default: "30s", OnEmpty: true}}
	out, results, err := envdefault.Apply(map[string]string{"TIMEOUT": ""}, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TIMEOUT"] != "30s" {
		t.Errorf("expected '30s', got %q", out["TIMEOUT"])
	}
	if len(results) != 1 || results[0].Reason != "empty" {
		t.Errorf("expected one empty result, got %+v", results)
	}
}

func TestApply_NilSecrets_ReturnsError(t *testing.T) {
	_, _, err := envdefault.Apply(nil, nil)
	if err == nil {
		t.Fatal("expected error for nil secrets")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"A": "1"}
	rules := []envdefault.Rule{{Key: "B", Default: "2"}}
	_, _, _ = envdefault.Apply(input, rules)
	if _, ok := input["B"]; ok {
		t.Error("input map was mutated")
	}
}

func TestLoadRules_EmptyPath(t *testing.T) {
	rules, err := envdefault.LoadRules("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected empty rules, got %d", len(rules))
	}
}

func TestLoadRules_Valid(t *testing.T) {
	content := `- key: DB_HOST
  default: localhost
- key: DB_PORT
  default: "5432"
  on_empty: true
`
	p := writeRulesFile(t, content)
	rules, err := envdefault.LoadRules(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[1].OnEmpty != true {
		t.Error("expected on_empty=true for second rule")
	}
}

func TestLoadRules_MissingFile(t *testing.T) {
	_, err := envdefault.LoadRules("/nonexistent/defaults.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
