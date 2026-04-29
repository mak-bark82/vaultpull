package envreplace

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a single replacement operation.
type Rule struct {
	Key     string // exact key to target, or empty to target all keys
	Pattern string // regex pattern to match in value
	With    string // replacement string (supports $1 capture groups)
}

// Result holds the outcome of a replacement operation.
type Result struct {
	Key      string
	OldValue string
	NewValue string
}

// Replace applies all rules to the provided env map and returns the
// modified map along with a list of changes made.
func Replace(env map[string]string, rules []Rule) (map[string]string, []Result, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var results []Result

	for _, rule := range rules {
		if rule.Pattern == "" {
			return nil, nil, fmt.Errorf("rule has empty pattern")
		}
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid pattern %q: %w", rule.Pattern, err)
		}

		for k, v := range out {
			if rule.Key != "" && rule.Key != k {
				continue
			}
			newVal := re.ReplaceAllString(v, rule.With)
			if newVal != v {
				results = append(results, Result{Key: k, OldValue: v, NewValue: newVal})
				out[k] = newVal
			}
		}
	}

	return out, results, nil
}

// Summary returns a human-readable summary of replacement results.
func Summary(results []Result) string {
	if len(results) == 0 {
		return "no replacements made"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d replacement(s) made:\n", len(results))
	for _, r := range results {
		fmt.Fprintf(&sb, "  %s: %q -> %q\n", r.Key, r.OldValue, r.NewValue)
	}
	return strings.TrimRight(sb.String(), "\n")
}
