package envrewrite

import (
	"fmt"
	"strings"
)

// Rule defines a single rewrite operation applied to env keys or values.
type Rule struct {
	// Key is the env variable name this rule targets. Empty means all keys.
	Key string `yaml:"key"`
	// Find is the substring or prefix to locate.
	Find string `yaml:"find"`
	// Replace is the replacement string.
	Replace string `yaml:"replace"`
	// Target specifies whether to rewrite "key", "value", or "both".
	Target string `yaml:"target"`
}

// Result holds the outcome of a rewrite operation.
type Result struct {
	Key     string
	OldKey  string
	OldVal  string
	NewVal  string
	Renamed bool
	Changed bool
}

// Rewrite applies the given rules to the secrets map and returns the
// rewritten map along with a slice of Result describing each change.
func Rewrite(secrets map[string]string, rules []Rule) (map[string]string, []Result) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var results []Result

	for _, rule := range rules {
		if rule.Find == "" {
			continue
		}
		target := strings.ToLower(rule.Target)
		if target == "" {
			target = "value"
		}

		keys := collectKeys(out, rule.Key)
		for _, k := range keys {
			v := out[k]
			res := Result{Key: k, OldKey: k, OldVal: v, NewVal: v}

			if target == "key" || target == "both" {
				newKey := strings.ReplaceAll(k, rule.Find, rule.Replace)
				if newKey != k {
					delete(out, k)
					out[newKey] = v
					res.Key = newKey
					res.Renamed = true
				}
			}

			if target == "value" || target == "both" {
				newVal := strings.ReplaceAll(v, rule.Find, rule.Replace)
				if newVal != v {
					out[res.Key] = newVal
					res.NewVal = newVal
					res.Changed = true
				}
			}

			if res.Renamed || res.Changed {
				results = append(results, res)
			}
		}
	}

	return out, results
}

// Summary returns a human-readable summary of rewrite results.
func Summary(results []Result) string {
	if len(results) == 0 {
		return "no rewrites applied"
	}
	var sb strings.Builder
	for _, r := range results {
		if r.Renamed && r.Changed {
			fmt.Fprintf(&sb, "renamed+changed: %s -> %s (value updated)\n", r.OldKey, r.Key)
		} else if r.Renamed {
			fmt.Fprintf(&sb, "renamed: %s -> %s\n", r.OldKey, r.Key)
		} else {
			fmt.Fprintf(&sb, "changed: %s value updated\n", r.Key)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func collectKeys(m map[string]string, filter string) []string {
	var keys []string
	for k := range m {
		if filter == "" || k == filter {
			keys = append(keys, k)
		}
	}
	return keys
}
