package envpurge

import (
	"testing"
)

var baseSecrets = map[string]string{
	"DB_PASSWORD": "secret",
	"API_KEY":     "key123",
	"APP_ENV":     "production",
	"LOG_LEVEL":   "info",
}

func TestPurge_RemovesMatchedKeys(t *testing.T) {
	rules := []Rule{{Key: "DB_PASSWORD"}, {Key: "API_KEY"}}
	res := Purge(baseSecrets, rules)
	if len(res.Removed) != 2 {
		t.Fatalf("expected 2 removed, got %d", len(res.Removed))
	}
	if _, ok := res.Removed["DB_PASSWORD"]; !ok {
		t.Error("expected DB_PASSWORD to be removed")
	}
}

func TestPurge_SkipsUnmatchedKeys(t *testing.T) {
	rules := []Rule{{Key: "DB_PASSWORD"}}
	res := Purge(baseSecrets, rules)
	if len(res.Skipped) != 3 {
		t.Fatalf("expected 3 skipped, got %d", len(res.Skipped))
	}
}

func TestPurge_NoRules_SkipsAll(t *testing.T) {
	res := Purge(baseSecrets, nil)
	if len(res.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(res.Removed))
	}
	if len(res.Skipped) != len(baseSecrets) {
		t.Errorf("expected all keys skipped")
	}
}

func TestApply_ReturnsRemainingKeys(t *testing.T) {
	rules := []Rule{{Key: "DB_PASSWORD"}, {Key: "API_KEY"}}
	out := Apply(baseSecrets, rules)
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been removed")
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %s", out["APP_ENV"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	copy := map[string]string{"SECRET": "val", "KEEP": "ok"}
	rules := []Rule{{Key: "SECRET"}}
	Apply(copy, rules)
	if _, ok := copy["SECRET"]; !ok {
		t.Error("Apply should not mutate the input map")
	}
}

func TestSummary_Format(t *testing.T) {
	res := Result{
		Removed: map[string]string{"A": "1", "B": "2"},
		Skipped: []string{"C"},
	}
	s := res.Summary()
	if s != "removed 2 key(s), skipped 1 key(s)" {
		t.Errorf("unexpected summary: %s", s)
	}
}
