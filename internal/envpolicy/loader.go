package envpolicy

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type rulesFile struct {
	Rules []Rule `yaml:"rules"`
}

// LoadRules reads a YAML file and returns the parsed policy rules.
// Returns an empty slice (no error) when path is empty.
func LoadRules(path string) ([]Rule, error) {
	if path == "" {
		return []Rule{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("envpolicy: read file: %w", err)
	}
	var rf rulesFile
	if err := yaml.Unmarshal(data, &rf); err != nil {
		return nil, fmt.Errorf("envpolicy: parse yaml: %w", err)
	}
	for i, r := range rf.Rules {
		if r.Name == "" {
			return nil, fmt.Errorf("envpolicy: rule at index %d is missing name", i)
		}
		if r.Target == "" {
			rf.Rules[i].Target = "key"
		}
		action := Action(strings.ToLower(string(r.Action)))
		if action != ActionAllow && action != ActionDeny && action != ActionWarn {
			return nil, fmt.Errorf("envpolicy: rule %q has unknown action %q", r.Name, r.Action)
		}
		rf.Rules[i].Action = action
	}
	return rf.Rules, nil
}
