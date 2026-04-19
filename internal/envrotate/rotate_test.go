package envrotate_test

import (
	"testing"

	"github.com/nicholasgasior/vaultpull/internal/envrotate"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_PASS": "old-pass",
		"API_KEY": "old-key",
		"STABLE":  "unchanged",
	}
}

func TestApply_DetectsChangedValues(t *testing.T) {
	r := envrotate.New(baseEnv())
	incoming := map[string]string{
		"DB_PASS": "new-pass",
		"API_KEY": "old-key", // unchanged
	}
	records, _ := r.Apply(incoming)
	if len(records) != 1 {
		t.Fatalf("expected 1 rotation record, got %d", len(records))
	}
	if records[0].Key != "DB_PASS" {
		t.Errorf("expected DB_PASS, got %s", records[0].Key)
	}
	if records[0].OldValue != "old-pass" || records[0].NewValue != "new-pass" {
		t.Errorf("unexpected old/new values: %+v", records[0])
	}
}

func TestApply_DetectsNewKey(t *testing.T) {
	r := envrotate.New(baseEnv())
	incoming := map[string]string{"NEW_SECRET": "value"}
	records, result := r.Apply(incoming)
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if result["NEW_SECRET"] != "value" {
		t.Errorf("new key not present in result")
	}
}

func TestApply_UnchangedProducesNoRecord(t *testing.T) {
	r := envrotate.New(baseEnv())
	incoming := map[string]string{"STABLE": "unchanged"}
	records, _ := r.Apply(incoming)
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}
}

func TestApply_ResultContainsAllKeys(t *testing.T) {
	r := envrotate.New(baseEnv())
	incoming := map[string]string{"DB_PASS": "rotated"}
	_, result := r.Apply(incoming)
	for _, k := range []string{"DB_PASS", "API_KEY", "STABLE"} {
		if _, ok := result[k]; !ok {
			t.Errorf("key %s missing from result", k)
		}
	}
}

func TestSummary_NoRecords(t *testing.T) {
	s := envrotate.Summary(nil)
	if s != "no secrets rotated" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_WithRecords(t *testing.T) {
	r := envrotate.New(baseEnv())
	incoming := map[string]string{"DB_PASS": "new", "API_KEY": "new-key"}
	records, _ := r.Apply(incoming)
	s := envrotate.Summary(records)
	if s != "2 secret(s) rotated" {
		t.Errorf("unexpected summary: %s", s)
	}
}
