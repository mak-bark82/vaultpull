package envlint_test

import (
	"testing"

	"github.com/nicholasgasior/vaultpull/internal/envlint"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"API_KEY":     "abc123",
		"SECRET_NAME": "my-secret",
	}
}

func TestLint_NoViolations(t *testing.T) {
	linter := envlint.New()
	violations := linter.Lint(baseSecrets())
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	linter := envlint.New()
	secrets := map[string]string{"DB_HOST": ""}
	violations := linter.Lint(secrets)
	if !hasRule(violations, "no-empty-value") {
		t.Error("expected no-empty-value violation")
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	linter := envlint.New()
	secrets := map[string]string{"db_host": "localhost"}
	violations := linter.Lint(secrets)
	if !hasRule(violations, "key-uppercase") {
		t.Error("expected key-uppercase violation")
	}
}

func TestLint_SpaceInKey(t *testing.T) {
	linter := envlint.New()
	secrets := map[string]string{"DB HOST": "localhost"}
	violations := linter.Lint(secrets)
	if !hasRule(violations, "no-spaces-in-key") {
		t.Error("expected no-spaces-in-key violation")
	}
}

func TestLint_InvalidKeyChars(t *testing.T) {
	linter := envlint.New()
	secrets := map[string]string{"DB-HOST": "localhost"}
	violations := linter.Lint(secrets)
	if !hasRule(violations, "valid-key-chars") {
		t.Error("expected valid-key-chars violation")
	}
}

func TestLint_MultipleViolations(t *testing.T) {
	linter := envlint.New()
	secrets := map[string]string{
		"db-host": "",
		"DB_PORT": "5432",
	}
	violations := linter.Lint(secrets)
	if len(violations) < 2 {
		t.Errorf("expected at least 2 violations, got %d", len(violations))
	}
}

func TestLint_WithCustomRules(t *testing.T) {
	rules := []envlint.Rule{
		{
			Name:    "no-localhost",
			Message: "value must not be localhost",
			Check:   func(_, value string) bool { return value == "localhost" },
		},
	}
	linter := envlint.WithRules(rules)
	secrets := map[string]string{"DB_HOST": "localhost"}
	violations := linter.Lint(secrets)
	if !hasRule(violations, "no-localhost") {
		t.Error("expected no-localhost violation from custom rule")
	}
}

func hasRule(violations []envlint.Violation, name string) bool {
	for _, v := range violations {
		if v.Rule == name {
			return true
		}
	}
	return false
}
