// Package envdrift detects configuration drift between a live Vault path
// and the current local .env file.
package envdrift

import "fmt"

// Status represents the drift state of a single key.
type Status int

const (
	StatusMatch   Status = iota // value is identical
	StatusChanged               // value differs
	StatusMissing               // key exists in Vault but not locally
	StatusExtra                 // key exists locally but not in Vault
)

// Entry describes the drift state for one environment key.
type Entry struct {
	Key        string
	Status     Status
	VaultValue string
	LocalValue string
}

// Report holds all drift entries for a single comparison.
type Report struct {
	Entries []Entry
}

// Detect compares vault secrets against local env values and returns a Report.
// vault and local are both map[string]string of key→value pairs.
func Detect(vault, local map[string]string) Report {
	seen := make(map[string]struct{})
	var entries []Entry

	for k, vv := range vault {
		seen[k] = struct{}{}
		lv, ok := local[k]
		switch {
		case !ok:
			entries = append(entries, Entry{Key: k, Status: StatusMissing, VaultValue: vv})
		case lv != vv:
			entries = append(entries, Entry{Key: k, Status: StatusChanged, VaultValue: vv, LocalValue: lv})
		default:
			entries = append(entries, Entry{Key: k, Status: StatusMatch, VaultValue: vv, LocalValue: lv})
		}
	}

	for k, lv := range local {
		if _, ok := seen[k]; !ok {
			entries = append(entries, Entry{Key: k, Status: StatusExtra, LocalValue: lv})
		}
	}

	return Report{Entries: entries}
}

// HasDrift returns true if any entry is not StatusMatch.
func (r Report) HasDrift() bool {
	for _, e := range r.Entries {
		if e.Status != StatusMatch {
			return true
		}
	}
	return false
}

// Summary returns a human-readable one-line summary of the report.
func (r Report) Summary() string {
	var changed, missing, extra int
	for _, e := range r.Entries {
		switch e.Status {
		case StatusChanged:
			changed++
		case StatusMissing:
			missing++
		case StatusExtra:
			extra++
		}
	}
	if !r.HasDrift() {
		return "no drift detected"
	}
	return fmt.Sprintf("drift detected: %d changed, %d missing, %d extra", changed, missing, extra)
}
