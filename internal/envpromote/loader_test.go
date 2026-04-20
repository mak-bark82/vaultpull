package envpromote_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/yourusername/vaultpull/internal/envpromote"
)

func writeRulesFile(t *testing.T, rs envpromote.RuleSet) string {
	t.Helper()
	f, err := os.CreateTemp("", "promote-rules-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(rs); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadRules_Valid(t *testing.T) {
	rs := envpromote.RuleSet{
		Rules: []envpromote.Rule{
			{Key: "API_KEY", FromEnv: "prod", ToEnv: "staging", Overwrite: true},
		},
	}
	path := writeRulesFile(t, rs)
	rules, err := envpromote.LoadRules(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || rules[0].Key != "API_KEY" {
		t.Error("expected one rule with key API_KEY")
	}
}

func TestLoadRules_EmptyPath(t *testing.T) {
	rules, err := envpromote.LoadRules("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Error("expected empty rules for empty path")
	}
}

func TestLoadRules_MissingFile(t *testing.T) {
	_, err := envpromote.LoadRules("/nonexistent/promote.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadRules_MissingKey(t *testing.T) {
	rs := envpromote.RuleSet{
		Rules: []envpromote.Rule{
			{Key: "", FromEnv: "prod", ToEnv: "staging"},
		},
	}
	path := writeRulesFile(t, rs)
	_, err := envpromote.LoadRules(path)
	if err == nil {
		t.Error("expected error for missing key")
	}
}
