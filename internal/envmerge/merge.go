package envmerge

import "maps"

// Strategy defines how conflicts are resolved when merging env maps.
type Strategy int

const (
	// PreferVault overwrites existing keys with values from Vault.
	PreferVault Strategy = iota
	// PreferLocal keeps existing local values on conflict.
	PreferLocal
)

// Result holds the merged map and metadata about the operation.
type Result struct {
	Merged    map[string]string
	Overwritten int
	Skipped     int
	Added       int
}

// Merge combines local and incoming (Vault) env maps according to the given strategy.
func Merge(local, incoming map[string]string, strategy Strategy) Result {
	result := Result{
		Merged: make(map[string]string),
	}

	// Seed with local values.
	maps.Copy(result.Merged, local)

	for k, v := range incoming {
		existing, exists := local[k]
		switch {
		case !exists:
			result.Merged[k] = v
			result.Added++
		case strategy == PreferVault && existing != v:
			result.Merged[k] = v
			result.Overwritten++
		case strategy == PreferLocal:
			result.Skipped++
		default:
			// Value unchanged — no action needed.
		}
	}

	return result
}
