package envdeprecate

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type rulesFile struct {
	Rules []Rule `yaml:"rules"`
}

// LoadRules reads deprecation rules from a YAML file.
// Returns an empty slice if path is empty.
func LoadRules(path string) ([]Rule, error) {
	if path == "" {
		return []Rule{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read rules file: %w", err)
	}
	var rf rulesFile
	if err := yaml.Unmarshal(data, &rf); err != nil {
		return nil, fmt.Errorf("parse rules file: %w", err)
	}
	for i, r := range rf.Rules {
		if r.Key == "" && r.Pattern == "" {
			return nil, fmt.Errorf("rule[%d]: must specify key or pattern", i)
		}
	}
	return rf.Rules, nil
}
