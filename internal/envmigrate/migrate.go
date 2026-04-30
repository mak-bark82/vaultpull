package envmigrate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a single migration step: rename a key, transform its value, or both.
type Rule struct {
	FromKey string `yaml:"from_key"`
	ToKey   string `yaml:"to_key"`
	Find    string `yaml:"find"`
	Replace string `yaml:"replace"`
}

// Result holds the outcome of applying a single migration rule.
type Result struct {
	Rule    Rule
	OldKey  string
	NewKey  string
	OldVal  string
	NewVal  string
	Applied bool
}

// Summary returns a human-readable summary of the migration results.
func Summary(results []Result) string {
	applied := 0
	for _, r := range results {
		if r.Applied {
			applied++
		}
	}
	return fmt.Sprintf("%d rule(s) applied out of %d", applied, len(results))
}

// Migrate applies the given rules to src, returning a new map and per-rule results.
// src is not mutated.
func Migrate(src map[string]string, rules []Rule) (map[string]string, []Result, error) {
	if src == nil {
		return nil, nil, fmt.Errorf("envmigrate: src must not be nil")
	}

	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	results := make([]Result, 0, len(rules))

	for _, rule := range rules {
		if rule.FromKey == "" {
			return nil, nil, fmt.Errorf("envmigrate: rule missing from_key")
		}

		val, exists := out[rule.FromKey]
		if !exists {
			results = append(results, Result{Rule: rule, OldKey: rule.FromKey, Applied: false})
			continue
		}

		res := Result{
			Rule:    rule,
			OldKey:  rule.FromKey,
			NewKey:  rule.FromKey,
			OldVal:  val,
			NewVal:  val,
			Applied: true,
		}

		// Apply value substitution if a find pattern is set.
		if rule.Find != "" {
			re, err := regexp.Compile(rule.Find)
			if err != nil {
				return nil, nil, fmt.Errorf("envmigrate: invalid find pattern %q: %w", rule.Find, err)
			}
			res.NewVal = re.ReplaceAllString(val, rule.Replace)
		}

		// Rename key if to_key is set and different.
		destKey := rule.FromKey
		if rule.ToKey != "" && rule.ToKey != rule.FromKey {
			destKey = strings.TrimSpace(rule.ToKey)
			delete(out, rule.FromKey)
			res.NewKey = destKey
		}

		out[destKey] = res.NewVal
		results = append(results, res)
	}

	return out, results, nil
}
