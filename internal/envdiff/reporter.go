package envdiff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ChangeKind represents the type of change detected between two env maps.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Changed  ChangeKind = "changed"
	Unchanged ChangeKind = "unchanged"
)

// Change describes a single key-level difference.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Report holds all detected changes between two env snapshots.
type Report struct {
	Changes []Change
}

// Diff compares oldEnv and newEnv and returns a Report of all differences.
func Diff(oldEnv, newEnv map[string]string) Report {
	seen := make(map[string]bool)
	var changes []Change

	for k, newVal := range newEnv {
		seen[k] = true
		if oldVal, exists := oldEnv[k]; !exists {
			changes = append(changes, Change{Key: k, Kind: Added, NewValue: newVal})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, Kind: Changed, OldValue: oldVal, NewValue: newVal})
		}
	}

	for k, oldVal := range oldEnv {
		if !seen[k] {
			changes = append(changes, Change{Key: k, Kind: Removed, OldValue: oldVal})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	return Report{Changes: changes}
}

// HasChanges returns true if the report contains any non-unchanged entries.
func (r Report) HasChanges() bool {
	return len(r.Changes) > 0
}

// Summary returns a brief count of added, removed, and changed keys.
func (r Report) Summary() string {
	var added, removed, changed int
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			added++
		case Removed:
			removed++
		case Changed:
			changed++
		}
	}
	return fmt.Sprintf("%d added, %d removed, %d changed", added, removed, changed)
}

// Print writes a human-readable summary of the report to w.
func (r Report) Print(w io.Writer) {
	if !r.HasChanges() {
		fmt.Fprintln(w, "No changes detected.")
		return
	}

	var sb strings.Builder
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			sb.WriteString(fmt.Sprintf("+ %s = %q\n", c.Key, c.NewValue))
		case Removed:
			sb.WriteString(fmt.Sprintf("- %s (was %q)\n", c.Key, c.OldValue))
		case Changed:
			sb.WriteString(fmt.Sprintf("~ %s: %q -> %q\n", c.Key, c.OldValue, c.NewValue))
		}
	}
	fmt.Fprint(w, sb.String())
}
