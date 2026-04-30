package envmigrate

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type rulesFile struct {
	Rules []Rule `yaml:"rules"`
}

// LoadRules reads migration rules from a YAML file at path.
// Returns an empty slice when path is empty.
func LoadRules(path string) ([]Rule, error) {
	if path == "" {
		return []Rule{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("envmigrate: reading rules file: %w", err)
	}

	var rf rulesFile
	if err := yaml.Unmarshal(data, &rf); err != nil {
		return nil, fmt.Errorf("envmigrate: parsing rules file: %w", err)
	}

	for i, r := range rf.Rules {
		if r.FromKey == "" {
			return nil, fmt.Errorf("envmigrate: rule[%d] missing required field 'from_key'", i)
		}
	}

	return rf.Rules, nil
}
