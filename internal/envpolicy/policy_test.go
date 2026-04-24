package envpolicy

import (
	"testing"
)

var baseSecrets = map[string]string{
	"DB_PASSWORD": "s3cr3t",
	"API_KEY":     "abc123",
	"DEBUG":       "true",
	"PLAIN":       "hello",
}

func TestNew_ValidRules(t *testing.T) {
	rules := []Rule{
		{Name: "no-debug", Pattern: `^DEBUG`, Target: "key", Action: ActionDeny},
	}
	e, err := New(rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil enforcer")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	rules := []Rule{
		{Name: "bad", Pattern: `[invalid`, Target: "key", Action: ActionDeny},
	}
	_, err := New(rules)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	rules := []Rule{
		{Name: "empty", Pattern: "", Target: "key", Action: ActionWarn},
	}
	_, err := New(rules)
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestCheck_DenyMatchesKey(t *testing.T) {
	rules := []Rule{
		{Name: "no-debug", Pattern: `^DEBUG$`, Target: "key", Action: ActionDeny},
	}
	e, _ := New(rules)
	vs := e.Check(baseSecrets)
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Action != ActionDeny {
		t.Errorf("expected deny, got %s", vs[0].Action)
	}
}

func TestCheck_WarnMatchesValue(t *testing.T) {
	rules := []Rule{
		{Name: "warn-true", Pattern: `^true$`, Target: "value", Action: ActionWarn},
	}
	e, _ := New(rules)
	vs := e.Check(baseSecrets)
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Key != "DEBUG" {
		t.Errorf("expected key DEBUG, got %s", vs[0].Key)
	}
}

func TestCheck_NoViolations(t *testing.T) {
	rules := []Rule{
		{Name: "no-xyz", Pattern: `^XYZ`, Target: "key", Action: ActionDeny},
	}
	e, _ := New(rules)
	vs := e.Check(baseSecrets)
	if len(vs) != 0 {
		t.Errorf("expected no violations, got %d", len(vs))
	}
}

func TestHasDenials_True(t *testing.T) {
	vs := []Violation{{Action: ActionWarn}, {Action: ActionDeny}}
	if !HasDenials(vs) {
		t.Error("expected HasDenials to return true")
	}
}

func TestHasDenials_False(t *testing.T) {
	vs := []Violation{{Action: ActionWarn}}
	if HasDenials(vs) {
		t.Error("expected HasDenials to return false")
	}
}
