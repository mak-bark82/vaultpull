package envaliases

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type aliasFile struct {
	Aliases []aliasEntry `yaml:"aliases"`
}

type aliasEntry struct {
	Name string   `yaml:"name"`
	Keys []string `yaml:"keys"`
}

// LoadAliases reads alias definitions from a YAML file at path.
// Returns an empty resolver when path is empty.
func LoadAliases(path string) (*Resolver, error) {
	if path == "" {
		return &Resolver{aliases: map[string][]string{}}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("envaliases: read file: %w", err)
	}

	var af aliasFile
	if err := yaml.Unmarshal(data, &af); err != nil {
		return nil, fmt.Errorf("envaliases: parse yaml: %w", err)
	}

	aliases := make([]Alias, 0, len(af.Aliases))
	for _, e := range af.Aliases {
		aliases = append(aliases, Alias{Name: e.Name, Keys: e.Keys})
	}

	return NewResolver(aliases)
}
