package envaliases

import (
	"fmt"
	"strings"
)

// Alias maps a short alias name to one or more environment variable keys.
type Alias struct {
	Name string
	Keys []string
}

// Resolver resolves alias names to their underlying env keys.
type Resolver struct {
	aliases map[string][]string
}

// NewResolver creates a Resolver from a slice of Alias definitions.
// Returns an error if any alias has an empty name or no keys.
func NewResolver(aliases []Alias) (*Resolver, error) {
	m := make(map[string][]string, len(aliases))
	for _, a := range aliases {
		if strings.TrimSpace(a.Name) == "" {
			return nil, fmt.Errorf("alias name must not be empty")
		}
		if len(a.Keys) == 0 {
			return nil, fmt.Errorf("alias %q must have at least one key", a.Name)
		}
		m[a.Name] = a.Keys
	}
	return &Resolver{aliases: m}, nil
}

// Resolve returns the env values for all keys mapped by the given alias.
// Unknown alias names return an error. Missing env keys are silently skipped.
func (r *Resolver) Resolve(alias string, env map[string]string) (map[string]string, error) {
	keys, ok := r.aliases[alias]
	if !ok {
		return nil, fmt.Errorf("unknown alias: %q", alias)
	}
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, exists := env[k]; exists {
			out[k] = v
		}
	}
	return out, nil
}

// Expand replaces alias placeholders in a key list with their underlying keys.
// Keys that are not aliases are passed through unchanged.
func (r *Resolver) Expand(keys []string) []string {
	var out []string
	for _, k := range keys {
		if mapped, ok := r.aliases[k]; ok {
			out = append(out, mapped...)
		} else {
			out = append(out, k)
		}
	}
	return out
}
