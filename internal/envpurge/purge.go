package envpurge

import "fmt"

// Rule defines a key pattern and reason for purging.
type Rule struct {
	Key    string `yaml:"key"`
	Reason string `yaml:"reason"`
}

// Result holds the outcome of a purge operation.
type Result struct {
	Removed map[string]string
	Skipped []string
}

// Summary returns a human-readable summary of the purge result.
func (r Result) Summary() string {
	return fmt.Sprintf("removed %d key(s), skipped %d key(s)", len(r.Removed), len(r.Skipped))
}

// Purge removes keys from secrets that match any of the provided rules.
// Keys not matched by any rule are collected in Skipped.
func Purge(secrets map[string]string, rules []Rule) Result {
	result := Result{
		Removed: make(map[string]string),
	}

	ruleIndex := make(map[string]struct{}, len(rules))
	for _, r := range rules {
		if r.Key != "" {
			ruleIndex[r.Key] = struct{}{}
		}
	}

	for k, v := range secrets {
		if _, matched := ruleIndex[k]; matched {
			result.Removed[k] = v
		} else {
			result.Skipped = append(result.Skipped, k)
		}
	}

	return result
}

// Apply returns a new map with matched keys removed.
func Apply(secrets map[string]string, rules []Rule) map[string]string {
	result := Purge(secrets, rules)
	out := make(map[string]string, len(result.Skipped))
	for _, k := range result.Skipped {
		out[k] = secrets[k]
	}
	return out
}
