// Package envreference detects and resolves cross-file secret references.
// A reference is expressed as ${file:KEY} inside a value, allowing one
// .env file to pull a resolved value from another already-loaded file.
package envreference

import (
	"fmt"
	"regexp"
	"strings"
)

// refPattern matches ${file:KEY} or ${KEY} style references.
var refPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Source holds a named set of key/value pairs that can be referenced.
type Source struct {
	Name   string
	Values map[string]string
}

// Resolver resolves cross-source references within env values.
type Resolver struct {
	sources map[string]map[string]string
}

// New creates a Resolver from a slice of named sources.
// Returns an error if any source name is empty.
func New(sources []Source) (*Resolver, error) {
	m := make(map[string]map[string]string, len(sources))
	for _, s := range sources {
		if strings.TrimSpace(s.Name) == "" {
			return nil, fmt.Errorf("envreference: source name must not be empty")
		}
		m[s.Name] = s.Values
	}
	return &Resolver{sources: m}, nil
}

// Resolve iterates over every value in env and expands any references it
// finds. References take the form ${sourceName:KEY}. If sourceName is
// omitted (plain ${KEY}), the resolver searches all sources in
// declaration order. Unresolvable references are left as-is.
func (r *Resolver) Resolve(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = r.expandValue(v)
	}
	return out
}

// expandValue replaces all reference tokens in a single value string.
func (r *Resolver) expandValue(v string) string {
	return refPattern.ReplaceAllStringFunc(v, func(match string) string {
		// Strip ${ and }
		inner := match[2 : len(match)-1]

		if idx := strings.Index(inner, ":"); idx >= 0 {
			// Explicit source reference: ${sourceName:KEY}
			srcName := inner[:idx]
			key := inner[idx+1:]
			if src, ok := r.sources[srcName]; ok {
				if val, ok := src[key]; ok {
					return val
				}
			}
			// Unresolvable — return original token unchanged.
			return match
		}

		// Implicit source: search all sources for the key.
		for _, src := range r.sources {
			if val, ok := src[inner]; ok {
				return val
			}
		}
		return match
	})
}

// FindReferences returns every unresolved reference token present in env
// values, deduplicated. Useful for validation and diagnostics.
func FindReferences(env map[string]string) []string {
	seen := make(map[string]struct{})
	var refs []string
	for _, v := range env {
		matches := refPattern.FindAllString(v, -1)
		for _, m := range matches {
			if _, ok := seen[m]; !ok {
				seen[m] = struct{}{}
				refs = append(refs, m)
			}
		}
	}
	return refs
}
