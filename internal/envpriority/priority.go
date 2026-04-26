package envpriority

import "fmt"

// Source represents a named environment source with an associated priority level.
// Lower numbers indicate higher priority (1 is highest).
type Source struct {
	Name     string
	Priority int
	Values   map[string]string
}

// Result holds the merged output and metadata about which source won each key.
type Result struct {
	Merged map[string]string
	Origin map[string]string // key -> source name
}

// Merge combines multiple sources according to their priority.
// If two sources define the same key, the one with the lower priority number wins.
// Sources with equal priority are resolved in the order they are provided.
func Merge(sources []Source) (Result, error) {
	for _, s := range sources {
		if s.Name == "" {
			return Result{}, fmt.Errorf("envpriority: source name must not be empty")
		}
		if s.Priority < 1 {
			return Result{}, fmt.Errorf("envpriority: source %q has invalid priority %d (must be >= 1)", s.Name, s.Priority)
		}
	}

	merged := make(map[string]string)
	origin := make(map[string]string)
	// Track the winning priority per key (lower = higher priority).
	winning := make(map[string]int)

	for _, src := range sources {
		for k, v := range src.Values {
			prev, seen := winning[k]
			if !seen || src.Priority < prev {
				merged[k] = v
				origin[k] = src.Name
				winning[k] = src.Priority
			}
		}
	}

	return Result{Merged: merged, Origin: origin}, nil
}

// Summary returns a human-readable description of how many keys each source contributed.
func Summary(r Result) map[string]int {
	counts := make(map[string]int)
	for _, src := range r.Origin {
		counts[src]++
	}
	return counts
}
