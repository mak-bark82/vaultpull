package envencrypt

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the YAML-decoded configuration for the encrypt processor.
type Config struct {
	Patterns []string `yaml:"patterns"`
}

// LoadConfig reads an envencrypt config YAML file from path.
// If path is empty, a zero-value Config is returned without error.
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return &Config{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("envencrypt: read config %q: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("envencrypt: parse config %q: %w", path, err)
	}
	if len(cfg.Patterns) == 0 {
		return nil, errors.New("envencrypt: config must define at least one pattern")
	}
	for i, p := range cfg.Patterns {
		if p == "" {
			return nil, fmt.Errorf("envencrypt: pattern at index %d must not be empty", i)
		}
	}
	return &cfg, nil
}
