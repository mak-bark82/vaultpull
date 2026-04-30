package envnormalize

import (
	"fmt"
	"regexp"
	"strings"
)

// Options controls normalization behaviour.
type Options struct {
	// UppercaseKeys converts all keys to UPPER_CASE.
	UppercaseKeys bool
	// ReplaceHyphen replaces hyphens in keys with underscores.
	ReplaceHyphen bool
	// StripInvalidChars removes characters from keys that are not
	// alphanumeric or underscores.
	StripInvalidChars bool
	// TrimValues trims leading/trailing whitespace from values.
	TrimValues bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys:     true,
		ReplaceHyphen:     true,
		StripInvalidChars: true,
		TrimValues:        true,
	}
}

var invalidKeyChars = regexp.MustCompile(`[^A-Za-z0-9_]`)

// Result holds the normalized map and a summary of changes.
type Result struct {
	Secrets  map[string]string
	Renamed  int
	Modified int
}

// Normalize applies the given options to secrets and returns a Result.
// Keys that collide after normalization return an error.
func Normalize(secrets map[string]string, opts Options) (Result, error) {
	out := make(map[string]string, len(secrets))
	var renamed, modified int

	for k, v := range secrets {
		newKey := k
		newVal := v

		if opts.ReplaceHyphen {
			newKey = strings.ReplaceAll(newKey, "-", "_")
		}
		if opts.StripInvalidChars {
			newKey = invalidKeyChars.ReplaceAllString(newKey, "")
		}
		if opts.UppercaseKeys {
			newKey = strings.ToUpper(newKey)
		}
		if opts.TrimValues {
			newVal = strings.TrimSpace(newVal)
		}

		if newKey == "" {
			continue
		}
		if _, exists := out[newKey]; exists {
			return Result{}, fmt.Errorf("key collision after normalization: %q", newKey)
		}

		if newKey != k {
			renamed++
		}
		if newVal != v {
			modified++
		}
		out[newKey] = newVal
	}

	return Result{Secrets: out, Renamed: renamed, Modified: modified}, nil
}
