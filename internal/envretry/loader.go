package envretry

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type policyFile struct {
	MaxAttempts int     `yaml:"max_attempts"`
	DelayMS     int64   `yaml:"delay_ms"`
	Multiplier  float64 `yaml:"multiplier"`
}

// LoadPolicy reads a YAML file and returns a Policy.
// If path is empty the DefaultPolicy is returned.
func LoadPolicy(path string) (Policy, error) {
	if path == "" {
		return DefaultPolicy(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Policy{}, fmt.Errorf("envretry: read policy file: %w", err)
	}

	var pf policyFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		return Policy{}, fmt.Errorf("envretry: parse policy file: %w", err)
	}

	if pf.MaxAttempts < 1 {
		return Policy{}, fmt.Errorf("envretry: max_attempts must be >= 1")
	}
	if pf.Multiplier <= 0 {
		pf.Multiplier = 1.0
	}

	return Policy{
		MaxAttempts: pf.MaxAttempts,
		Delay:       time.Duration(pf.DelayMS) * time.Millisecond,
		Multiplier:  pf.Multiplier,
	}, nil
}
