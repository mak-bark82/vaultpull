package envcleanup_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/envcleanup"
)

var baseEnv = map[string]string{
	"DATABASE_URL": "  postgres://localhost/db  ",
	"API_KEY":      `"my-secret-key"`,
	"EMPTY_VAR":    "   ",
	"PLAIN_VAR":    "hello",
}

func TestApply_TrimSpace(t *testing.T) {
	rules := []envcleanup.Rule{
		{Key: "DATABASE_URL", TrimSpace: true},
	}
	out, _, _ := envcleanup.Apply(baseEnv, rules)
	if got := out["DATABASE_URL"]; got != "postgres://localhost/db" {
		t.Errorf("expected trimmed value, got %q", got)
	}
}

func TestApply_StripQuotes(t *testing.T) {
	rules := []envcleanup.Rule{
		{Key: "API_KEY", StripQuotes: true},
	}
	out, _, _ := envcleanup.Apply(baseEnv, rules)
	if got := out["API_KEY"]; got != "my-secret-key" {
		t.Errorf("expected unquoted value, got %q", got)
	}
}

func TestApply_RemoveEmpty(t *testing.T) {
	rules := []envcleanup.Rule{
		{Key: "EMPTY_VAR", TrimSpace: true, RemoveEmpty: true},
	}
	out, results, summary := envcleanup.Apply(baseEnv, rules)
	if _, exists := out["EMPTY_VAR"]; exists {
		t.Error("expected EMPTY_VAR to be removed")
	}
	if summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", summary.Removed)
	}
	var found bool
	for _, r := range results {
		if r.Key == "EMPTY_VAR" && r.Removed {
			found = true
		}
	}
	if !found {
		t.Error("expected result entry with Removed=true for EMPTY_VAR")
	}
}

func TestApply_GlobPattern(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "  localhost  ",
		"DB_PORT": "  5432  ",
		"APP_ENV": "  production  ",
	}
	rules := []envcleanup.Rule{
		{Key: "DB_*", TrimSpace: true},
	}
	out, _, summary := envcleanup.Apply(env, rules)
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST trimmed, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT trimmed, got %q", out["DB_PORT"])
	}
	if out["APP_ENV"] != "  production  " {
		t.Errorf("expected APP_ENV unchanged, got %q", out["APP_ENV"])
	}
	if summary.Changed != 2 {
		t.Errorf("expected 2 changed, got %d", summary.Changed)
	}
}

func TestApply_NoRules_NoChange(t *testing.T) {
	out, _, summary := envcleanup.Apply(baseEnv, nil)
	if summary.Changed != 0 || summary.Removed != 0 {
		t.Errorf("expected no changes, got changed=%d removed=%d", summary.Changed, summary.Removed)
	}
	if len(out) != len(baseEnv) {
		t.Errorf("expected same number of keys")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"KEY": "  value  "}
	rules := []envcleanup.Rule{{Key: "KEY", TrimSpace: true}}
	envcleanup.Apply(env, rules)
	if env["KEY"] != "  value  " {
		t.Error("Apply must not mutate the input map")
	}
}
