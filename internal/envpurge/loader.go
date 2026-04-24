package envpurge

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type rulesFile struct {
	Rules []Rule `yaml:"rules"`
}

// LoadRules reads purge rules from a YAML file at the given path.
// Returns an empty slice if path is empty.
func LoadRules(path string) ([]Rule, error) {
	if path == "" {
		return []Rule{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("envpurge: read rules file: %w", err)
	}

	var rf rulesFile
	if err := yaml.Unmarshal(data, &rf); err != nil {
		return nil, fmt.Errorf("envpurge: parse rules file: %w", err)
	}

	for i, r := range rf.Rules {
		if r.Key == "" {
			return nil, fmt.Errorf("envpurge: rule %d is missing 'key'", i)
		}
	}

	return rf.Rules, nil
}
