package envquota

import (
	"testing"
)

var baseSecrets = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "s3cr3t",
	"API_KEY":     "abc123",
}

func TestCheck_NoViolations(t *testing.T) {
	rule := Rule{MaxKeys: 10, MaxKeyLength: 50, MaxValLength: 100}
	res, err := Check(baseSecrets, rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.OK() {
		t.Errorf("expected no violations, got %d", len(res.Violations))
	}
}

func TestCheck_TooManyKeys(t *testing.T) {
	rule := Rule{MaxKeys: 2}
	res, err := Check(baseSecrets, rule)
	if err == nil {
		t.Fatal("expected error for too many keys")
	}
	if res.OK() {
		t.Error("expected violations")
	}
	found := false
	for _, v := range res.Violations {
		if v.Key == "__total__" {
			found = true
		}
	}
	if !found {
		t.Error("expected __total__ violation")
	}
}

func TestCheck_KeyTooLong(t *testing.T) {
	secrets := map[string]string{"VERY_LONG_KEY_NAME": "value"}
	rule := Rule{MaxKeyLength: 5}
	res, err := Check(secrets, rule)
	if err == nil {
		t.Fatal("expected error for long key")
	}
	if len(res.Violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(res.Violations))
	}
}

func TestCheck_ValueTooLong(t *testing.T) {
	secrets := map[string]string{"KEY": "this-value-is-way-too-long"}
	rule := Rule{MaxValLength: 5}
	res, err := Check(secrets, rule)
	if err == nil {
		t.Fatal("expected error for long value")
	}
	if res.Violations[0].Key != "KEY" {
		t.Errorf("expected violation for KEY, got %s", res.Violations[0].Key)
	}
}

func TestCheck_ZeroLimitsAreUnlimited(t *testing.T) {
	rule := Rule{} // all zeros = unlimited
	_, err := Check(baseSecrets, rule)
	if err != nil {
		t.Fatalf("unexpected error with zero limits: %v", err)
	}
}

func TestSummary_PassAndFail(t *testing.T) {
	pass := Result{}
	if pass.Summary() != "quota check passed: no violations" {
		t.Errorf("unexpected summary: %s", pass.Summary())
	}

	fail := Result{Violations: []Violation{{Key: "X", Message: "too long"}}}
	if fail.Summary() != "quota check failed: 1 violation(s)" {
		t.Errorf("unexpected summary: %s", fail.Summary())
	}
}
