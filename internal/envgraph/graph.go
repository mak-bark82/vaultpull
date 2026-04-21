package envgraph

import (
	"fmt"
	"sort"
	"strings"
)

// Node represents a single env variable and its dependencies.
type Node struct {
	Key  string
	Deps []string
}

// Graph holds a dependency graph of env variables.
type Graph struct {
	nodes map[string]*Node
}

// New creates an empty Graph.
func New() *Graph {
	return &Graph{nodes: make(map[string]*Node)}
}

// Add registers a key with its list of dependency keys.
func (g *Graph) Add(key string, deps []string) {
	g.nodes[key] = &Node{Key: key, Deps: deps}
}

// Resolve returns keys in topological order (dependencies before dependents).
// Returns an error if a cycle is detected.
func (g *Graph) Resolve() ([]string, error) {
	visited := make(map[string]bool)
	onStack := make(map[string]bool)
	var order []string

	var visit func(key string) error
	visit = func(key string) error {
		if onStack[key] {
			return fmt.Errorf("cycle detected at key: %s", key)
		}
		if visited[key] {
			return nil
		}
		onStack[key] = true
		if node, ok := g.nodes[key]; ok {
			for _, dep := range node.Deps {
				if err := visit(dep); err != nil {
					return err
				}
			}
		}
		onStack[key] = false
		visited[key] = true
		order = append(order, key)
		return nil
	}

	keys := make([]string, 0, len(g.nodes))
	for k := range g.nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if err := visit(k); err != nil {
			return nil, err
		}
	}
	return order, nil
}

// BuildFromEnv constructs a Graph by scanning env values for ${VAR} references.
func BuildFromEnv(env map[string]string) *Graph {
	g := New()
	for key, val := range env {
		deps := extractRefs(val)
		g.Add(key, deps)
	}
	return g
}

// extractRefs returns all ${KEY} references found in a value string.
func extractRefs(val string) []string {
	var refs []string
	for {
		start := strings.Index(val, "${") 
		if start == -1 {
			break
		}
		end := strings.Index(val[start:], "}")
		if end == -1 {
			break
		}
		ref := val[start+2 : start+end]
		if ref != "" {
			refs = append(refs, ref)
		}
		val = val[start+end+1:]
	}
	return refs
}
