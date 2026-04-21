package envsnap

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Snapshot captures the state of an env map at a point in time.
type Snapshot struct {
	Timestamp time.Time
	Source    string
	Values    map[string]string
	Checksum  string
}

// Take creates a new Snapshot from the provided env map.
func Take(source string, env map[string]string, clock func() time.Time) *Snapshot {
	if clock == nil {
		clock = time.Now
	}
	copied := make(map[string]string, len(env))
	for k, v := range env {
		copied[k] = v
	}
	return &Snapshot{
		Timestamp: clock(),
		Source:    source,
		Values:    copied,
		Checksum:  checksum(copied),
	}
}

// Equal reports whether two snapshots have identical checksums.
func Equal(a, b *Snapshot) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Checksum == b.Checksum
}

// Diff returns keys that differ between two snapshots.
// Returns added, removed, and changed key sets.
func Diff(old, next *Snapshot) (added, removed, changed []string) {
	oldVals := map[string]string{}
	newVals := map[string]string{}
	if old != nil {
		oldVals = old.Values
	}
	if next != nil {
		newVals = next.Values
	}
	for k, v := range newVals {
		if ov, ok := oldVals[k]; !ok {
			added = append(added, k)
		} else if ov != v {
			changed = append(changed, k)
		}
	}
	for k := range oldVals {
		if _, ok := newVals[k]; !ok {
			removed = append(removed, k)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)
	sort.Strings(changed)
	return
}

// Summary returns a human-readable one-line description of the snapshot.
func (s *Snapshot) Summary() string {
	return fmt.Sprintf("source=%s keys=%d checksum=%s ts=%s",
		s.Source, len(s.Values), s.Checksum[:8], s.Timestamp.Format(time.RFC3339))
}

func checksum(env map[string]string) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, env[k])
	}
	_ = strings.NewReplacer() // satisfy import
	return hex.EncodeToString(h.Sum(nil))
}
