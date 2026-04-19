package envrotate

import (
	"fmt"
	"time"
)

// RotationRecord holds metadata about a secret rotation event.
type RotationRecord struct {
	Key       string
	OldValue  string
	NewValue  string
	RotatedAt time.Time
}

// Rotator applies new secret values and tracks what changed.
type Rotator struct {
	current map[string]string
}

// New creates a new Rotator seeded with the current env state.
func New(current map[string]string) *Rotator {
	copy := make(map[string]string, len(current))
	for k, v := range current {
		copy[k] = v
	}
	return &Rotator{current: copy}
}

// Apply merges incoming secrets into the current state, returning rotation
// records for every key whose value actually changed.
func (r *Rotator) Apply(incoming map[string]string) ([]RotationRecord, map[string]string) {
	var records []RotationRecord
	now := time.Now().UTC()

	result := make(map[string]string, len(r.current))
	for k, v := range r.current {
		result[k] = v
	}

	for k, newVal := range incoming {
		oldVal, exists := r.current[k]
		if !exists || oldVal != newVal {
			records = append(records, RotationRecord{
				Key:       k,
				OldValue:  oldVal,
				NewValue:  newVal,
				RotatedAt: now,
			})
			result[k] = newVal
		}
	}

	return records, result
}

// Summary returns a human-readable rotation summary.
func Summary(records []RotationRecord) string {
	if len(records) == 0 {
		return "no secrets rotated"
	}
	return fmt.Sprintf("%d secret(s) rotated", len(records))
}
