package envquota

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type ruleFile struct {
	MaxKeys      int `yaml:"max_keys"`
	MaxKeyLength int `yaml:"max_key_length"`
	MaxValLength int `yaml:"max_val_length"`
}

// LoadRule reads a YAML file and returns a Rule.
// If path is empty, a zero-value Rule (no limits) is returned.
func LoadRule(path string) (Rule, error) {
	if path == "" {
		return Rule{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Rule{}, err
	}

	var rf ruleFile
	if err := yaml.Unmarshal(data, &rf); err != nil {
		return Rule{}, err
	}

	if rf.MaxKeys < 0 || rf.MaxKeyLength < 0 || rf.MaxValLength < 0 {
		return Rule{}, errors.New("quota: limits must be non-negative")
	}

	return Rule{
		MaxKeys:      rf.MaxKeys,
		MaxKeyLength: rf.MaxKeyLength,
		MaxValLength: rf.MaxValLength,
	}, nil
}
