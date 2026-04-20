package envreport

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// Entry represents a single secret key's sync status in a report.
type Entry struct {
	Key    string
	Status string // "added", "changed", "removed", "unchanged"
	Source string // vault path
}

// Report holds a summary of a sync operation.
type Report struct {
	Timestamp time.Time
	EnvFile   string
	Entries   []Entry
}

// New creates a new Report for the given env file.
func New(envFile string, entries []Entry) *Report {
	return &Report{
		Timestamp: time.Now().UTC(),
		EnvFile:   envFile,
		Entries:   entries,
	}
}

// Summary returns counts of each status type.
func (r *Report) Summary() map[string]int {
	counts := map[string]int{
		"added":     0,
		"changed":   0,
		"removed":   0,
		"unchanged": 0,
	}
	for _, e := range r.Entries {
		if _, ok := counts[e.Status]; ok {
			counts[e.Status]++
		}
	}
	return counts
}

// Render writes a human-readable report to w.
func (r *Report) Render(w io.Writer) error {
	fmt.Fprintf(w, "vaultpull sync report\n")
	fmt.Fprintf(w, "Timestamp : %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "Env file  : %s\n", r.EnvFile)
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 48))

	sorted := make([]Entry, len(r.Entries))
	copy(sorted, r.Entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, e := range sorted {
		fmt.Fprintf(w, "  [%-9s] %s  (%s)\n", e.Status, e.Key, e.Source)
	}

	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 48))
	s := r.Summary()
	fmt.Fprintf(w, "added: %d  changed: %d  removed: %d  unchanged: %d\n",
		s["added"], s["changed"], s["removed"], s["unchanged"])
	return nil
}
