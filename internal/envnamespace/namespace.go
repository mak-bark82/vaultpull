package envnamespace

import (
	"fmt"
	"strings"
)

// Namespace represents a named prefix scope for environment variable keys.
type Namespace struct {
	Name   string
	Prefix string
}

// Resolver maps keys into and out of namespaces.
type Resolver struct {
	namespaces map[string]Namespace
}

// NewResolver creates a Resolver from a slice of Namespace definitions.
// Returns an error if any namespace has an empty name or prefix.
func NewResolver(ns []Namespace) (*Resolver, error) {
	m := make(map[string]Namespace, len(ns))
	for _, n := range ns {
		if n.Name == "" {
			return nil, fmt.Errorf("namespace name must not be empty")
		}
		if n.Prefix == "" {
			return nil, fmt.Errorf("namespace %q: prefix must not be empty", n.Name)
		}
		m[n.Name] = n
	}
	return &Resolver{namespaces: m}, nil
}

// Qualify adds the namespace prefix to every key in the map.
// Returns an error if the namespace name is unknown.
func (r *Resolver) Qualify(name string, env map[string]string) (map[string]string, error) {
	ns, ok := r.namespaces[name]
	if !ok {
		return nil, fmt.Errorf("unknown namespace %q", name)
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[ns.Prefix+k] = v
	}
	return out, nil
}

// Strip removes the namespace prefix from keys that carry it.
// Keys without the prefix are passed through unchanged.
func (r *Resolver) Strip(name string, env map[string]string) (map[string]string, error) {
	ns, ok := r.namespaces[name]
	if !ok {
		return nil, fmt.Errorf("unknown namespace %q", name)
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		if strings.HasPrefix(k, ns.Prefix) {
			out[strings.TrimPrefix(k, ns.Prefix)] = v
		} else {
			out[k] = v
		}
	}
	return out, nil
}

// Names returns all registered namespace names.
func (r *Resolver) Names() []string {
	names := make([]string, 0, len(r.namespaces))
	for n := range r.namespaces {
		names = append(names, n)
	}
	return names
}
