package envprofile

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Profile represents a named environment profile (e.g., dev, staging, prod)
type Profile struct {
	Name        string            `yaml:"name"`
	VaultPrefix string            `yaml:"vault_prefix"`
	EnvFile     string            `yaml:"env_file"`
	Overrides   map[string]string `yaml:"overrides,omitempty"`
}

// ProfileSet holds multiple named profiles loaded from a profiles config file.
type ProfileSet struct {
	Profiles map[string]Profile `yaml:"profiles"`
}

// LoadProfiles reads a YAML file and returns a ProfileSet.
func LoadProfiles(path string) (*ProfileSet, error) {
	if path == "" {
		return nil, fmt.Errorf("profile path must not be empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading profile file: %w", err)
	}

	var ps ProfileSet
	if err := yaml.Unmarshal(data, &ps); err != nil {
		return nil, fmt.Errorf("parsing profile file: %w", err)
	}

	for name, p := range ps.Profiles {
		if p.VaultPrefix == "" {
			return nil, fmt.Errorf("profile %q missing vault_prefix", name)
		}
		if p.EnvFile == "" {
			return nil, fmt.Errorf("profile %q missing env_file", name)
		}
	}

	return &ps, nil
}

// Get returns a profile by name or an error if not found.
func (ps *ProfileSet) Get(name string) (Profile, error) {
	p, ok := ps.Profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("profile %q not found", name)
	}
	return p, nil
}
