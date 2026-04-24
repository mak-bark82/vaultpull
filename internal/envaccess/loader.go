package envaccess

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ruleFile struct {
	Rules []struct {
		Pattern    string `yaml:"pattern"`
		Permission string `yaml:"permission"`
	} `yaml:"rules"`
}

// LoadRules reads access rules from a YAML file.
// Returns an empty slice when path is empty.
func LoadRules(path string) ([]Rule, error) {
	if path == "" {
		return []Rule{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("envaccess: read file: %w", err)
	}

	var rf ruleFile
	if err := yaml.Unmarshal(data, &rf); err != nil {
		return nil, fmt.Errorf("envaccess: parse yaml: %w", err)
	}

	var rules []Rule
	for _, r := range rf.Rules {
		perm, err := parsePermission(r.Permission)
		if err != nil {
			return nil, fmt.Errorf("envaccess: rule %q: %w", r.Pattern, err)
		}
		rules = append(rules, Rule{Pattern: r.Pattern, Permission: perm})
	}
	return rules, nil
}

func parsePermission(s string) (Permission, error) {
	switch s {
	case "read":
		return PermRead, nil
	case "write":
		return PermWrite, nil
	case "none", "":
		return PermNone, nil
	default:
		return PermNone, errors.New("unknown permission: " + s)
	}
}
