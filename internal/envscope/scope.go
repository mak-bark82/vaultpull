package envscope

import "fmt"

// Scope defines a named environment scope (e.g. dev, staging, prod)
// and maps it to a Vault path prefix.
type Scope struct {
	Name   string `yaml:"name"`
	Prefix string `yaml:"prefix"`
}

// Resolver resolves secret paths for a given scope.
type Resolver struct {
	scopes map[string]Scope
}

// NewResolver creates a Resolver from a slice of Scopes.
func NewResolver(scopes []Scope) (*Resolver, error) {
	m := make(map[string]Scope, len(scopes))
	for _, s := range scopes {
		if s.Name == "" {
			return nil, fmt.Errorf("scope missing name")
		}
		if s.Prefix == "" {
			return nil, fmt.Errorf("scope %q missing prefix", s.Name)
		}
		m[s.Name] = s
	}
	return &Resolver{scopes: m}, nil
}

// Resolve returns the full Vault path for the given scope and relative path.
func (r *Resolver) Resolve(scopeName, path string) (string, error) {
	s, ok := r.scopes[scopeName]
	if !ok {
		return "", fmt.Errorf("unknown scope %q", scopeName)
	}
	return s.Prefix + "/" + path, nil
}

// Names returns all registered scope names.
func (r *Resolver) Names() []string {
	names := make([]string, 0, len(r.scopes))
	for k := range r.scopes {
		names = append(names, k)
	}
	return names
}
