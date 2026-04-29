package envlease

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds loader configuration for a lease manager.
type Config struct {
	Path        string        `yaml:"path"`
	DefaultTTL  time.Duration `yaml:"default_ttl"`
}

// LoadConfig reads a YAML config file for the lease manager.
// Returns a default Config if path is empty.
func LoadConfig(path string) (Config, error) {
	if path == "" {
		return Config{DefaultTTL: time.Hour}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, errors.New("lease config file not found: " + path)
		}
		return Config{}, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	if cfg.DefaultTTL <= 0 {
		cfg.DefaultTTL = time.Hour
	}
	return cfg, nil
}
