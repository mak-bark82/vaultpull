// Package envdefault applies default values to missing or empty keys
// in a secret map based on a set of declarative rules.
package envdefault

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Rule defines a default value for a specific key.
type Rule struct {
	Key     string `yaml:"key"`
	Default string `yaml:"default"`
	OnEmpty bool   `yaml:"on_empty"` // also apply when key exists but value is empty
}

// Result holds the outcome of applying defaults.
type Result struct {
	Key    string
	Value  string
	Reason string // "missing" or "empty"
}

// Apply applies the given rules to secrets, filling in defaults where needed.
// It returns the updated map and a slice of results describing changes made.
func Apply(secrets map[string]string, rules []Rule) (map[string]string, []Result, error) {
	if secrets == nil {
		return nil, nil, errors.New("envdefault: secrets map must not be nil")
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var results []Result
	for _, r := range rules {
		if r.Key == "" {
			continue
		}
		v, exists := out[r.Key]
		switch {
		case !exists:
			out[r.Key] = r.Default
			results = append(results, Result{Key: r.Key, Value: r.Default, Reason: "missing"})
		case r.OnEmpty && v == "":
			out[r.Key] = r.Default
			results = append(results, Result{Key: r.Key, Value: r.Default, Reason: "empty"})
		}
	}
	return out, results, nil
}

// LoadRules reads a YAML file containing a list of default rules.
// Returns an empty slice if path is empty.
func LoadRules(path string) ([]Rule, error) {
	if path == "" {
		return []Rule{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("envdefault: read rules: %w", err)
	}
	var rules []Rule
	if err := yaml.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("envdefault: parse rules: %w", err)
	}
	return rules, nil
}
