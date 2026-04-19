package envvalidator_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/envvalidator"
)

func TestValidate_AllValid(t *testing.T) {
	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	rules := []envvalidator.Rule{
		{Key: "APP_ENV", Required: true, Allowed: []string{"production", "staging", "development"}},
		{Key: "PORT", Required: true},
	}
	if err := envvalidator.Validate(env, rules); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	env := map[string]string{}
	rules := []envvalidator.Rule{
		{Key: "DATABASE_URL", Required: true},
	}
	err := envvalidator.Validate(env, rules)
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
	ve, ok := err.(*envvalidator.ValidationError)
	if !ok || len(ve.Errors) != 1 {
		t.Fatalf("expected 1 validation error, got: %v", err)
	}
}

func TestValidate_EmptyRequired(t *testing.T) {
	env := map[string]string{"SECRET": "   "}
	rules := []envvalidator.Rule{
		{Key: "SECRET", Required: true},
	}
	if err := envvalidator.Validate(env, rules); err == nil {
		t.Fatal("expected error for empty required value")
	}
}

func TestValidate_InvalidAllowedValue(t *testing.T) {
	env := map[string]string{"LOG_LEVEL": "verbose"}
	rules := []envvalidator.Rule{
		{Key: "LOG_LEVEL", Allowed: []string{"debug", "info", "warn", "error"}},
	}
	err := envvalidator.Validate(env, rules)
	if err == nil {
		t.Fatal("expected error for disallowed value")
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	env := map[string]string{"APP_ENV": "unknown"}
	rules := []envvalidator.Rule{
		{Key: "APP_ENV", Required: true, Allowed: []string{"production", "staging"}},
		{Key: "PORT", Required: true},
	}
	err := envvalidator.Validate(env, rules)
	ve, ok := err.(*envvalidator.ValidationError)
	if !ok || len(ve.Errors) != 2 {
		t.Fatalf("expected 2 validation errors, got: %v", err)
	}
}

func TestValidate_OptionalMissingKey(t *testing.T) {
	env := map[string]string{}
	rules := []envvalidator.Rule{
		{Key: "OPTIONAL_KEY", Required: false, Allowed: []string{"yes", "no"}},
	}
	if err := envvalidator.Validate(env, rules); err != nil {
		t.Fatalf("expected no error for missing optional key, got: %v", err)
	}
}
