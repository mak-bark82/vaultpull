package envreplace_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envreplace"
)

var baseEnv = map[string]string{
	"DB_HOST": "localhost",
	"DB_PORT": "5432",
	"API_URL": "http://old.example.com/api",
	"SECRET":  "top-secret-value",
}

func TestReplace_AllKeys(t *testing.T) {
	rules := []envreplace.Rule{
		{Pattern: "localhost", With: "db.prod.internal"},
	}
	out, results, err := envreplace.Replace(baseEnv, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "db.prod.internal" {
		t.Errorf("expected DB_HOST to be replaced, got %q", out["DB_HOST"])
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestReplace_TargetedKey(t *testing.T) {
	rules := []envreplace.Rule{
		{Key: "API_URL", Pattern: "http", With: "https"},
	}
	out, results, err := envreplace.Replace(baseEnv, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_URL"] != "https://old.example.com/api" {
		t.Errorf("unexpected API_URL: %q", out["API_URL"])
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	// DB_HOST should be unchanged
	if out["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST should be unchanged, got %q", out["DB_HOST"])
	}
}

func TestReplace_NoMatch(t *testing.T) {
	rules := []envreplace.Rule{
		{Pattern: "nonexistent", With: "x"},
	}
	_, results, err := envreplace.Replace(baseEnv, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestReplace_InvalidPattern(t *testing.T) {
	rules := []envreplace.Rule{
		{Pattern: "[", With: "x"},
	}
	_, _, err := envreplace.Replace(baseEnv, rules)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestReplace_EmptyPattern_ReturnsError(t *testing.T) {
	rules := []envreplace.Rule{
		{Pattern: "", With: "x"},
	}
	_, _, err := envreplace.Replace(baseEnv, rules)
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestReplace_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	rules := []envreplace.Rule{{Pattern: "value", With: "replaced"}}
	_, _, _ = envreplace.Replace(env, rules)
	if env["KEY"] != "value" {
		t.Error("input map was mutated")
	}
}

func TestReplace_CaptureGroup(t *testing.T) {
	env := map[string]string{"ENDPOINT": "v1/users"}
	rules := []envreplace.Rule{
		{Pattern: `v(\d+)/`, With: "v${1}-stable/"},
	}
	out, _, err := envreplace.Replace(env, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ENDPOINT"] != "v1-stable/users" {
		t.Errorf("unexpected value: %q", out["ENDPOINT"])
	}
}

func TestSummary_NoResults(t *testing.T) {
	s := envreplace.Summary(nil)
	if s != "no replacements made" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestSummary_WithResults(t *testing.T) {
	results := []envreplace.Result{
		{Key: "FOO", OldValue: "bar", NewValue: "baz"},
	}
	s := envreplace.Summary(results)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
