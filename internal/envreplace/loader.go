package envreplace

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type rulesFile struct {
	Rules []struct {
		Key     string `yaml:"key"`
		Pattern string `yaml:"pattern"`
		With    string `yaml:"with"`
	} `yaml:"rules"`
}

// LoadRules reads replacement rules from a YAML file.
// Returns an empty slice if path is empty.
func LoadRules(path string) ([]Rule, error) {
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read rules file: %w", err)
	}

	var rf rulesFile
	if err := yaml.Unmarshal(data, &rf); err != nil {
		return nil, fmt.Errorf("parse rules file: %w", err)
	}

	rules := make([]Rule, 0, len(rf.Rules))
	for i, r := range rf.Rules {
		if r.Pattern == "" {
			return nil, fmt.Errorf("rule[%d]: missing pattern", i)
		}
		rules = append(rules, Rule{
			Key:     r.Key,
			Pattern: r.Pattern,
			With:    r.With,
		})
	}
	return rules, nil
}
