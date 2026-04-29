// Package envtrim provides utilities for trimming and normalizing
// environment variable keys and values according to configurable rules.
package envtrim

import (
	"strings"
)

// Options controls which trimming operations are applied.
type Options struct {
	// TrimKeys removes leading/trailing whitespace from keys.
	TrimKeys bool
	// TrimValues removes leading/trailing whitespace from values.
	TrimValues bool
	// StripValueQuotes removes surrounding single or double quotes from values.
	StripValueQuotes bool
	// NormalizeKeys converts keys to uppercase.
	NormalizeKeys bool
}

// DefaultOptions returns an Options with safe, common defaults.
func DefaultOptions() Options {
	return Options{
		TrimKeys:         true,
		TrimValues:       true,
		StripValueQuotes: false,
		NormalizeKeys:    false,
	}
}

// Trim applies the given Options to every entry in env and returns a new map.
// The original map is never mutated.
func Trim(env map[string]string, opts Options) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		if opts.TrimKeys {
			k = strings.TrimSpace(k)
		}
		if opts.NormalizeKeys {
			k = strings.ToUpper(k)
		}
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		if opts.StripValueQuotes {
			v = stripQuotes(v)
		}
		if k != "" {
			result[k] = v
		}
	}
	return result
}

// TrimValue applies trimming rules to a single value string.
func TrimValue(v string, opts Options) string {
	if opts.TrimValues {
		v = strings.TrimSpace(v)
	}
	if opts.StripValueQuotes {
		v = stripQuotes(v)
	}
	return v
}

func stripQuotes(s string) string {
	if len(s) < 2 {
		return s
	}
	if (s[0] == '"' && s[len(s)-1] == '"') ||
		(s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}
