package envpromote

import "fmt"

// Rule defines how a key should be promoted from one environment to another.
type Rule struct {
	Key       string
	FromEnv   string
	ToEnv     string
	Overwrite bool
}

// Result captures what happened to a single key during promotion.
type Result struct {
	Key     string
	Skipped bool
	Reason  string
}

// Promote copies keys from src into dst according to the given rules.
// Returns a slice of Results describing each action taken.
func Promote(src, dst map[string]string, rules []Rule) (map[string]string, []Result) {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	var results []Result

	for _, rule := range rules {
		val, ok := src[rule.Key]
		if !ok {
			results = append(results, Result{
				Key:     rule.Key,
				Skipped: true,
				Reason:  fmt.Sprintf("key %q not found in source (%s)", rule.Key, rule.FromEnv),
			})
			continue
		}

		if existing, exists := out[rule.Key]; exists && !rule.Overwrite {
			results = append(results, Result{
				Key:     rule.Key,
				Skipped: true,
				Reason:  fmt.Sprintf("key %q already exists in target (%s) with value %q", rule.Key, rule.ToEnv, existing),
			})
			continue
		}

		out[rule.Key] = val
		results = append(results, Result{
			Key:     rule.Key,
			Skipped: false,
			Reason:  fmt.Sprintf("promoted from %s to %s", rule.FromEnv, rule.ToEnv),
		})
	}

	return out, results
}
