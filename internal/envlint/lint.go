package envlint

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a single lint rule applied to env keys or values.
type Rule struct {
	Name    string
	Message string
	Check   func(key, value string) bool
}

// Violation represents a single lint failure.
type Violation struct {
	Key     string
	Rule    string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("[%s] %s: %s", v.Rule, v.Key, v.Message)
}

var defaultRules = []Rule{
	{
		Name:    "no-empty-value",
		Message: "value must not be empty",
		Check:   func(_, value string) bool { return strings.TrimSpace(value) == "" },
	},
	{
		Name:    "key-uppercase",
		Message: "key must be uppercase",
		Check:   func(key, _ string) bool { return key != strings.ToUpper(key) },
	},
	{
		Name:    "no-spaces-in-key",
		Message: "key must not contain spaces",
		Check:   func(key, _ string) bool { return strings.Contains(key, " ") },
	},
	{
		Name:    "valid-key-chars",
		Message: "key must contain only letters, digits, and underscores",
		Check: func(key, _ string) bool {
			matched, _ := regexp.MatchString(`^[A-Z0-9_]+$`, key)
			return !matched
		},
	},
}

// Linter runs lint rules against a set of env secrets.
type Linter struct {
	rules []Rule
}

// New creates a Linter with the default rule set.
func New() *Linter {
	return &Linter{rules: defaultRules}
}

// WithRules creates a Linter with a custom rule set.
func WithRules(rules []Rule) *Linter {
	return &Linter{rules: rules}
}

// Lint checks all secrets against the configured rules and returns any violations.
func (l *Linter) Lint(secrets map[string]string) []Violation {
	var violations []Violation
	for key, value := range secrets {
		for _, rule := range l.rules {
			if rule.Check(key, value) {
				violations = append(violations, Violation{
					Key:     key,
					Rule:    rule.Name,
					Message: rule.Message,
				})
			}
		}
	}
	return violations
}
