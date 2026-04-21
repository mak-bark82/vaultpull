package envgraph

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ExportDOT writes the dependency graph in Graphviz DOT format to w.
func (g *Graph) ExportDOT(w io.Writer) error {
	_, err := fmt.Fprintln(w, "digraph envgraph {")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, `  rankdir=LR;`)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(g.nodes))
	for k := range g.nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		node := g.nodes[key]
		if len(node.Deps) == 0 {
			_, err = fmt.Fprintf(w, "  %q;\n", sanitize(key))
			if err != nil {
				return err
			}
			continue
		}
		for _, dep := range node.Deps {
			_, err = fmt.Fprintf(w, "  %q -> %q;\n", sanitize(dep), sanitize(key))
			if err != nil {
				return err
			}
		}
	}

	_, err = fmt.Fprintln(w, "}")
	return err
}

// sanitize removes characters unsafe for DOT identifiers.
func sanitize(s string) string {
	return strings.ReplaceAll(s, "\"", "_")
}
