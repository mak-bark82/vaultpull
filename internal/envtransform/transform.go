package envtransform

import (
	"strings"
)

// Rule defines a transformation to apply to secret values.
type Rule struct {
	Prefix    string
	Suffix    string
	UpperCase bool
	LowerCase bool
	TrimSpace bool
}

// Transform applies the given Rule to a map of env key/value pairs
// and returns a new map with transformed values.
func Transform(secrets map[string]string, rule Rule) map[string]string {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		v = apply(v, rule)
		result[k] = v
	}
	return result
}

// TransformValue applies a Rule to a single value.
func TransformValue(v string, rule Rule) string {
	return apply(v, rule)
}

func apply(v string, rule Rule) string {
	if rule.TrimSpace {
		v = strings.TrimSpace(v)
	}
	if rule.UpperCase {
		v = strings.ToUpper(v)
	} else if rule.LowerCase {
		v = strings.ToLower(v)
	}
	if rule.Prefix != "" {
		v = rule.Prefix + v
	}
	if rule.Suffix != "" {
		v = v + rule.Suffix
	}
	return v
}
