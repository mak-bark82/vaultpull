package envrewrite

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// RuleSet is a named collection of rewrite rules loaded from YAML.
type RuleSet struct {
	Rules []Rule `yaml:"rules"`
}

// LoadRules reads a YAML file at path and returns the parsed RuleSet.
// If path is empty, an empty RuleSet is returned without error.
func LoadRules(path string) (*RuleSet, error) {
	if path == "" {
		return &RuleSet{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("envrewrite: read rules file: %w", err)
	}

	var rs RuleSet
	if err := yaml.Unmarshal(data, &rs); err != nil {
		return nil, fmt.Errorf("envrewrite: parse rules file: %w", err)
	}

	for i, r := range rs.Rules {
		if r.Find == "" {
			return nil, fmt.Errorf("envrewrite: rule[%d] missing 'find' field", i)
		}
		target := r.Target
		if target != "" && target != "key" && target != "value" && target != "both" {
			return nil, fmt.Errorf("envrewrite: rule[%d] invalid target %q (must be key, value, or both)", i, target)
		}
	}

	return &rs, nil
}
