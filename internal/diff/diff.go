package diff

import "fmt"

// Result holds the comparison between existing and incoming secrets.
type Result struct {
	Added   map[string]string
	Changed map[string]string
	Removed map[string]string
	Unchanged map[string]string
}

// Compare computes the diff between existing env vars and incoming secrets.
func Compare(existing, incoming map[string]string) Result {
	r := Result{
		Added:     make(map[string]string),
		Changed:   make(map[string]string),
		Removed:   make(map[string]string),
		Unchanged: make(map[string]string),
	}

	for k, v := range incoming {
		if old, ok := existing[k]; !ok {
			r.Added[k] = v
		} else if old != v {
			r.Changed[k] = v
		} else {
			r.Unchanged[k] = v
		}
	}

	for k, v := range existing {
		if _, ok := incoming[k]; !ok {
			r.Removed[k] = v
		}
	}

	return r
}

// Summary returns a human-readable summary of the diff.
func (r Result) Summary() string {
	return fmt.Sprintf("added=%d changed=%d removed=%d unchanged=%d",
		len(r.Added), len(r.Changed), len(r.Removed), len(r.Unchanged))
}

// HasChanges returns true if there is anything to apply.
func (r Result) HasChanges() bool {
	return len(r.Added) > 0 || len(r.Changed) > 0 || len(r.Removed) > 0
}
