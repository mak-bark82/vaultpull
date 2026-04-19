package envvalidator

import (
	"fmt"
	"strings"
)

// Rule defines a validation rule for an environment variable.
type Rule struct {
	Key      string
	Required bool
	Allowed  []string // if non-empty, value must be one of these
}

// ValidationError holds all validation failures.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return "validation failed:\n  " + strings.Join(e.Errors, "\n  ")
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

// Validate checks the provided env map against the given rules.
// Returns a *ValidationError if any rules are violated, nil otherwise.
func Validate(env map[string]string, rules []Rule) error {
	ve := &ValidationError{}

	for _, rule := range rules {
		val, exists := env[rule.Key]

		if rule.Required && (!exists || strings.TrimSpace(val) == "") {
			ve.Errors = append(ve.Errors, fmt.Sprintf("required key %q is missing or empty", rule.Key))
			continue
		}

		if exists && len(rule.Allowed) > 0 {
			if !contains(rule.Allowed, val) {
				ve.Errors = append(ve.Errors, fmt.Sprintf(
					"key %q has invalid value %q; allowed values: %s",
					rule.Key, val, strings.Join(rule.Allowed, ", "),
				))
			}
		}
	}

	if ve.HasErrors() {
		return ve
	}
	return nil
}

func contains(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}
