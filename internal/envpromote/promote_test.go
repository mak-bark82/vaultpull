package envpromote_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/envpromote"
)

func baseSrc() map[string]string {
	return map[string]string{
		"DB_HOST": "prod.db.example.com",
		"API_KEY": "prod-secret-key",
	}
}

func baseDst() map[string]string {
	return map[string]string{
		"DB_HOST": "staging.db.example.com",
	}
}

func TestPromote_CopiesKey(t *testing.T) {
	rules := []envpromote.Rule{
		{Key: "API_KEY", FromEnv: "prod", ToEnv: "staging", Overwrite: false},
	}
	out, results := envpromote.Promote(baseSrc(), baseDst(), rules)
	if out["API_KEY"] != "prod-secret-key" {
		t.Errorf("expected API_KEY to be promoted, got %q", out["API_KEY"])
	}
	if len(results) != 1 || results[0].Skipped {
		t.Error("expected one non-skipped result")
	}
}

func TestPromote_SkipsExistingWithoutOverwrite(t *testing.T) {
	rules := []envpromote.Rule{
		{Key: "DB_HOST", FromEnv: "prod", ToEnv: "staging", Overwrite: false},
	}
	out, results := envpromote.Promote(baseSrc(), baseDst(), rules)
	if out["DB_HOST"] != "staging.db.example.com" {
		t.Error("expected existing value to be preserved")
	}
	if !results[0].Skipped {
		t.Error("expected result to be skipped")
	}
}

func TestPromote_OverwritesWhenAllowed(t *testing.T) {
	rules := []envpromote.Rule{
		{Key: "DB_HOST", FromEnv: "prod", ToEnv: "staging", Overwrite: true},
	}
	out, results := envpromote.Promote(baseSrc(), baseDst(), rules)
	if out["DB_HOST"] != "prod.db.example.com" {
		t.Error("expected value to be overwritten")
	}
	if results[0].Skipped {
		t.Error("expected result not to be skipped")
	}
}

func TestPromote_SkipsMissingSourceKey(t *testing.T) {
	rules := []envpromote.Rule{
		{Key: "MISSING_KEY", FromEnv: "prod", ToEnv: "staging", Overwrite: false},
	}
	_, results := envpromote.Promote(baseSrc(), baseDst(), rules)
	if !results[0].Skipped {
		t.Error("expected missing key to be skipped")
	}
}

func TestPromote_DoesNotMutateInput(t *testing.T) {
	src := baseSrc()
	dst := baseDst()
	rules := []envpromote.Rule{
		{Key: "API_KEY", FromEnv: "prod", ToEnv: "staging", Overwrite: true},
	}
	envpromote.Promote(src, dst, rules)
	if _, ok := dst["API_KEY"]; ok {
		t.Error("expected dst to not be mutated")
	}
}
