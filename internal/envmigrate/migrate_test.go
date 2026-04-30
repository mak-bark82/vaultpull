package envmigrate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/subtlepseudonym/vaultpull/internal/envmigrate"
)

func writeRulesFile(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "rules.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

var baseSecrets = map[string]string{
	"OLD_DB_HOST": "localhost",
	"OLD_DB_PORT": "5432",
	"API_SECRET":  "hunter2",
}

func TestMigrate_RenamesKey(t *testing.T) {
	rules := []envmigrate.Rule{{FromKey: "OLD_DB_HOST", ToKey: "DB_HOST"}}
	out, results, err := envmigrate.Migrate(baseSecrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["OLD_DB_HOST"]; ok {
		t.Error("expected OLD_DB_HOST to be removed")
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if !results[0].Applied {
		t.Error("expected result to be applied")
	}
}

func TestMigrate_TransformsValue(t *testing.T) {
	rules := []envmigrate.Rule{{FromKey: "OLD_DB_PORT", Find: `5432`, Replace: "5433"}}
	out, _, err := envmigrate.Migrate(baseSecrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["OLD_DB_PORT"] != "5433" {
		t.Errorf("expected OLD_DB_PORT=5433, got %q", out["OLD_DB_PORT"])
	}
}

func TestMigrate_MissingKeyNotApplied(t *testing.T) {
	rules := []envmigrate.Rule{{FromKey: "DOES_NOT_EXIST", ToKey: "NEW_KEY"}}
	_, results, err := envmigrate.Migrate(baseSecrets, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Applied {
		t.Error("expected result not to be applied for missing key")
	}
}

func TestMigrate_DoesNotMutateInput(t *testing.T) {
	orig := map[string]string{"OLD_KEY": "value"}
	rules := []envmigrate.Rule{{FromKey: "OLD_KEY", ToKey: "NEW_KEY"}}
	_, _, err := envmigrate.Migrate(orig, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := orig["OLD_KEY"]; !ok {
		t.Error("Migrate must not mutate input map")
	}
}

func TestMigrate_NilSrcReturnsError(t *testing.T) {
	_, _, err := envmigrate.Migrate(nil, nil)
	if err == nil {
		t.Error("expected error for nil src")
	}
}

func TestSummary(t *testing.T) {
	results := []envmigrate.Result{
		{Applied: true},
		{Applied: false},
		{Applied: true},
	}
	s := envmigrate.Summary(results)
	if s != "2 rule(s) applied out of 3" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestLoadRules_Valid(t *testing.T) {
	p := writeRulesFile(t, "rules:\n  - from_key: OLD_KEY\n    to_key: NEW_KEY\n")
	rules, err := envmigrate.LoadRules(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || rules[0].FromKey != "OLD_KEY" {
		t.Errorf("unexpected rules: %+v", rules)
	}
}

func TestLoadRules_EmptyPath(t *testing.T) {
	rules, err := envmigrate.LoadRules("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected empty rules, got %d", len(rules))
	}
}

func TestLoadRules_MissingFromKey(t *testing.T) {
	p := writeRulesFile(t, "rules:\n  - to_key: NEW_KEY\n")
	_, err := envmigrate.LoadRules(p)
	if err == nil {
		t.Error("expected error for missing from_key")
	}
}
