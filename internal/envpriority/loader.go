package envpriority

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// SourceConfig represents a single priority source entry in the YAML config file.
type SourceConfig struct {
	Name     string            `yaml:"name"`
	Priority int               `yaml:"priority"`
	Values   map[string]string `yaml:"values"`
}

// FileConfig is the top-level structure of a priority sources YAML file.
type FileConfig struct {
	Sources []SourceConfig `yaml:"sources"`
}

// LoadSources reads a YAML file defining priority-ordered env sources and
// returns them as a slice of Source ready for use with Merge.
//
// The YAML format is:
//
//	 sources:
//	   - name: vault
//	     priority: 10
//	     values:
//	       DB_HOST: vault-db
//	   - name: local
//	     priority: 5
//	     values:
//	       DB_HOST: local-db
//	       DB_PORT: "5432"
func LoadSources(path string) ([]Source, error) {
	if path == "" {
		return nil, fmt.Errorf("envpriority: sources file path must not be empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("envpriority: failed to read sources file %q: %w", path, err)
	}

	var fc FileConfig
	if err := yaml.Unmarshal(data, &fc); err != nil {
		return nil, fmt.Errorf("envpriority: failed to parse sources file %q: %w", path, err)
	}

	if len(fc.Sources) == 0 {
		return nil, fmt.Errorf("envpriority: sources file %q contains no sources", path)
	}

	sources := make([]Source, 0, len(fc.Sources))
	for i, sc := range fc.Sources {
		if sc.Name == "" {
			return nil, fmt.Errorf("envpriority: source at index %d is missing a name", i)
		}
		if sc.Priority < 0 {
			return nil, fmt.Errorf("envpriority: source %q has negative priority %d", sc.Name, sc.Priority)
		}
		sources = append(sources, Source{
			Name:     sc.Name,
			Priority: sc.Priority,
			Values:   sc.Values,
		})
	}

	return sources, nil
}
