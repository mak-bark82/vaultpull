package envscope

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type scopeFile struct {
	Scopes []Scope `yaml:"scopes"`
}

// LoadScopes reads a YAML file and returns a Resolver.
// If path is empty, returns an empty Resolver.
func LoadScopes(path string) (*Resolver, error) {
	if path == "" {
		return &Resolver{scopes: map[string]Scope{}}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read scopes file: %w", err)
	}
	var sf scopeFile
	if err := yaml.Unmarshal(data, &sf); err != nil {
		return nil, fmt.Errorf("parse scopes file: %w", err)
	}
	return NewResolver(sf.Scopes)
}
