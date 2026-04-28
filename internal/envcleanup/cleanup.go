package envcleanup

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a cleanup rule applied to env values.
type Rule struct {
	Key     string // exact key or glob pattern (e.g. "*_URL")
	TrimSpace   bool
	StripQuotes bool
	RemoveEmpty bool
}

// Result holds the outcome of applying cleanup to a single key.
type Result struct {
	Key     string
	OldValue string
	NewValue string
	Changed  bool
	Removed  bool
}

// Summary describes the overall cleanup operation.
type Summary struct {
	Changed int
	Removed int
	Total   int
}

// Apply runs all cleanup rules against the provided env map.
// It returns the cleaned map, a list of per-key results, and a summary.
func Apply(env map[string]string, rules []Rule) (map[string]string, []Result, Summary) {
	out := make(map[string]string, len(env))
	var results []Result
	var summary Summary

	for k, v := range env {
		summary.Total++
		newVal := v
		removed := false

		for _, r := range rules {
			if !matchesKey(r.Key, k) {
				continue
			}
			if r.TrimSpace {
				newVal = strings.TrimSpace(newVal)
			}
			if r.StripQuotes {
				newVal = strings.Trim(newVal, `"'`)
			}
			if r.RemoveEmpty && strings.TrimSpace(newVal) == "" {
				removed = true
			}
		}

		res := Result{Key: k, OldValue: v, NewValue: newVal}
		if removed {
			res.Removed = true
			summary.Removed++
		} else {
			out[k] = newVal
			if newVal != v {
				res.Changed = true
				summary.Changed++
			}
		}
		results = append(results, res)
	}

	return out, results, summary
}

// matchesKey returns true if pattern matches key.
// Supports exact match and simple glob with leading/trailing *.
func matchesKey(pattern, key string) bool {
	if pattern == "*" {
		return true
	}
	if !strings.ContainsAny(pattern, "*?") {
		return pattern == key
	}
	regexPat := "^" + regexp.QuoteMeta(pattern) + "$"
	regexPat = strings.ReplaceAll(regexPat, `\*`, ".*")
	regexPat = strings.ReplaceAll(regexPat, `\?`, ".")
	matched, err := regexp.MatchString(regexPat, key)
	if err != nil {
		fmt.Printf("envcleanup: invalid pattern %q: %v\n", pattern, err)
		return false
	}
	return matched
}
