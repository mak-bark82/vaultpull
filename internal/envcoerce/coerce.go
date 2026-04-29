// Package envcoerce provides type coercion utilities for environment variable values.
// It converts string values to target types and normalises common representations.
package envcoerce

import (
	"fmt"
	"strconv"
	"strings"
)

// TargetType represents the desired output type for a coercion rule.
type TargetType string

const (
	TypeBool   TargetType = "bool"
	TypeInt    TargetType = "int"
	TypeFloat  TargetType = "float"
	TypeString TargetType = "string"
)

// Rule maps an environment variable key to its desired type.
type Rule struct {
	Key  string
	Type TargetType
}

// Result holds the outcome of a coercion attempt.
type Result struct {
	Key      string
	Original string
	Coerced  string
	Changed  bool
	Err      error
}

// Coerce applies the provided rules to the given env map, returning a new map
// with values normalised to their canonical string representation, along with
// a slice of per-key results describing what changed or failed.
func Coerce(env map[string]string, rules []Rule) (map[string]string, []Result) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	results := make([]Result, 0, len(rules))
	for _, rule := range rules {
		orig, ok := out[rule.Key]
		if !ok {
			continue
		}
		coerced, err := coerceValue(orig, rule.Type)
		r := Result{
			Key:      rule.Key,
			Original: orig,
			Coerced:  coerced,
			Changed:  coerced != orig,
			Err:      err,
		}
		if err == nil {
			out[rule.Key] = coerced
		}
		results = append(results, r)
	}
	return out, results
}

// coerceValue normalises a raw string value to the canonical form of the
// requested target type.
func coerceValue(raw string, t TargetType) (string, error) {
	switch t {
	case TypeBool:
		b, err := strconv.ParseBool(strings.TrimSpace(raw))
		if err != nil {
			return raw, fmt.Errorf("cannot coerce %q to bool: %w", raw, err)
		}
		return strconv.FormatBool(b), nil
	case TypeInt:
		i, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err != nil {
			return raw, fmt.Errorf("cannot coerce %q to int: %w", raw, err)
		}
		return strconv.FormatInt(i, 10), nil
	case TypeFloat:
		f, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			return raw, fmt.Errorf("cannot coerce %q to float: %w", raw, err)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case TypeString:
		return strings.TrimSpace(raw), nil
	default:
		return raw, fmt.Errorf("unknown target type %q", t)
	}
}
