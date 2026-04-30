package envclassify

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type rulesFile struct {
	Rules []Rule `yaml:"rules"`
}

// LoadRules reads classification rules from a YAML file.
// Returns an empty slice when path is empty.
func LoadRules(path string) ([]Rule, error) {
	if path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("envclassify: read rules file: %w", err)
	}
	var rf rulesFile
	if err := yaml.Unmarshal(data, &rf); err != nil {
		return nil, fmt.Errorf("envclassify: parse rules file: %w", err)
	}
	for i, r := range rf.Rules {
		if r.Pattern == "" {
			return nil, fmt.Errorf("envclassify: rule[%d] missing pattern", i)
		}
		if r.Category == "" {
			return nil, fmt.Errorf("envclassify: rule[%d] missing category", i)
		}
	}
	return rf.Rules, nil
}
