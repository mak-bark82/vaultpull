package envquota

import (
	"errors"
	"fmt"
)

// Rule defines a quota constraint for env secrets.
type Rule struct {
	MaxKeys      int // 0 means unlimited
	MaxKeyLength int // 0 means unlimited
	MaxValLength int // 0 means unlimited
}

// Violation describes a single quota breach.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("quota violation [%s]: %s", v.Key, v.Message)
}

// Result holds all violations found during a Check.
type Result struct {
	Violations []Violation
}

// OK returns true when no violations were found.
func (r Result) OK() bool { return len(r.Violations) == 0 }

// Summary returns a human-readable summary string.
func (r Result) Summary() string {
	if r.OK() {
		return "quota check passed: no violations"
	}
	return fmt.Sprintf("quota check failed: %d violation(s)", len(r.Violations))
}

// Check validates env secrets against the given Rule.
// It returns a Result and a non-nil error when any violation is found.
func Check(secrets map[string]string, rule Rule) (Result, error) {
	var result Result

	if rule.MaxKeys > 0 && len(secrets) > rule.MaxKeys {
		result.Violations = append(result.Violations, Violation{
			Key:     "__total__",
			Message: fmt.Sprintf("key count %d exceeds max %d", len(secrets), rule.MaxKeys),
		})
	}

	for k, v := range secrets {
		if rule.MaxKeyLength > 0 && len(k) > rule.MaxKeyLength {
			result.Violations = append(result.Violations, Violation{
				Key:     k,
				Message: fmt.Sprintf("key length %d exceeds max %d", len(k), rule.MaxKeyLength),
			})
		}
		if rule.MaxValLength > 0 && len(v) > rule.MaxValLength {
			result.Violations = append(result.Violations, Violation{
				Key:     k,
				Message: fmt.Sprintf("value length %d exceeds max %d", len(v), rule.MaxValLength),
			})
		}
	}

	if !result.OK() {
		return result, errors.New(result.Summary())
	}
	return result, nil
}
