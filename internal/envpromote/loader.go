package envpromote

import (
	"encoding/json"
	"fmt"
	"os"
)

// RuleSet is the top-level structure for a promote rules JSON file.
type RuleSet struct {
	Rules []Rule `json:"rules"`
}

// LoadRules reads promotion rules from a JSON file at the given path.
// Returns an empty slice if path is empty.
func LoadRules(path string) ([]Rule, error) {
	if path == "" {
		return []Rule{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("envpromote: reading rules file: %w", err)
	}

	var rs RuleSet
	if err := json.Unmarshal(data, &rs); err != nil {
		return nil, fmt.Errorf("envpromote: parsing rules file: %w", err)
	}

	for i, r := range rs.Rules {
		if r.Key == "" {
			return nil, fmt.Errorf("envpromote: rule[%d] missing key", i)
		}
		if r.FromEnv == "" || r.ToEnv == "" {
			return nil, fmt.Errorf("envpromote: rule[%d] missing from_env or to_env", i)
		}
	}

	return rs.Rules, nil
}
