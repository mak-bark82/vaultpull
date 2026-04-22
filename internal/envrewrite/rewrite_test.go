package envrewrite

import (
	"strings"
	"testing"
)

var baseSecrets = map[string]string{
	"APP_HOST":    "localhost",
	"APP_PORT":    "8080",
	"DB_HOST":     "db.internal",
	"DB_PASSWORD": "secret123",
}

func TestRewrite_ReplaceValue(t *testing.T) {
	rules := []Rule{
		{Key: "APP_HOST", Find: "localhost", Replace: "prod.example.com", Target: "value"},
	}
	out, results := Rewrite(baseSecrets, rules)
	if out["APP_HOST"] != "prod.example.com" {
		t.Errorf("expected prod.example.com, got %s", out["APP_HOST"])
	}
	if len(results) != 1 || !results[0].Changed {
		t.Error("expected one changed result")
	}
}

func TestRewrite_RenameKey(t *testing.T) {
	rules := []Rule{
		{Key: "APP_PORT", Find: "APP_", Replace: "SVC_", Target: "key"},
	}
	out, results := Rewrite(baseSecrets, rules)
	if _, ok := out["SVC_PORT"]; !ok {
		t.Error("expected SVC_PORT to exist after rename")
	}
	if _, ok := out["APP_PORT"]; ok {
		t.Error("expected APP_PORT to be removed after rename")
	}
	if len(results) != 1 || !results[0].Renamed {
		t.Error("expected one renamed result")
	}
}

func TestRewrite_Both_KeyAndValue(t *testing.T) {
	rules := []Rule{
		{Key: "DB_HOST", Find: "DB_", Replace: "DATABASE_", Target: "both"},
	}
	secrets := map[string]string{"DB_HOST": "DB_server"}
	out, results := Rewrite(secrets, rules)
	if _, ok := out["DATABASE_HOST"]; !ok {
		t.Error("expected DATABASE_HOST key")
	}
	if out["DATABASE_HOST"] != "DATABASE_server" {
		t.Errorf("unexpected value: %s", out["DATABASE_HOST"])
	}
	if len(results) != 1 || !results[0].Renamed || !results[0].Changed {
		t.Error("expected renamed and changed result")
	}
}

func TestRewrite_NoMatchProducesNoResult(t *testing.T) {
	rules := []Rule{
		{Key: "APP_HOST", Find: "nonexistent", Replace: "x", Target: "value"},
	}
	out, results := Rewrite(baseSecrets, rules)
	if out["APP_HOST"] != "localhost" {
		t.Error("value should be unchanged")
	}
	if len(results) != 0 {
		t.Error("expected no results")
	}
}

func TestRewrite_EmptyFindSkipsRule(t *testing.T) {
	rules := []Rule{
		{Key: "APP_HOST", Find: "", Replace: "x", Target: "value"},
	}
	_, results := Rewrite(baseSecrets, rules)
	if len(results) != 0 {
		t.Error("empty Find should skip rule")
	}
}

func TestRewrite_DoesNotMutateInput(t *testing.T) {
	orig := map[string]string{"KEY": "original"}
	rules := []Rule{{Key: "KEY", Find: "original", Replace: "changed", Target: "value"}}
	Rewrite(orig, rules)
	if orig["KEY"] != "original" {
		t.Error("input map should not be mutated")
	}
}

func TestSummary_NoResults(t *testing.T) {
	s := Summary(nil)
	if s != "no rewrites applied" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_ContainsKeyInfo(t *testing.T) {
	results := []Result{
		{Key: "NEW_KEY", OldKey: "OLD_KEY", Renamed: true},
	}
	s := Summary(results)
	if !strings.Contains(s, "OLD_KEY") || !strings.Contains(s, "NEW_KEY") {
		t.Errorf("summary missing key info: %s", s)
	}
}
