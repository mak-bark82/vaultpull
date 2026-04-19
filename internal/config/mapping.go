package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Mapping defines a single Vault path -> local env file relationship.
type Mapping struct {
	VaultPath string `json:"vault_path"`
	EnvFile   string `json:"env_file"`
	Overwrite bool   `json:"overwrite"`
}

// LoadMappings reads a JSON file containing a list of Mapping entries.
// If the file path is empty, it returns an empty slice without error.
func LoadMappings(path string) ([]Mapping, error) {
	if path == "" {
		return []Mapping{}, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("mappings: cannot open %s: %w", path, err)
	}
	defer f.Close()

	var mappings []Mapping
	if err := json.NewDecoder(f).Decode(&mappings); err != nil {
		return nil, fmt.Errorf("mappings: invalid JSON in %s: %w", path, err)
	}

	for i, m := range mappings {
		if m.VaultPath == "" {
			return nil, fmt.Errorf("mappings[%d]: vault_path is required", i)
		}
		if m.EnvFile == "" {
			return nil, fmt.Errorf("mappings[%d]: env_file is required", i)
		}
	}

	return mappings, nil
}
