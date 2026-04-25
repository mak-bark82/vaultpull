package envfreeze

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// FrozenSet represents a set of env keys that are locked against modification.
type FrozenSet struct {
	Keys      []string  `json:"keys"`
	FrozenAt  time.Time `json:"frozen_at"`
	Comment   string    `json:"comment,omitempty"`
}

// Freezer manages a frozen key set persisted to disk.
type Freezer struct {
	path string
	set  *FrozenSet
}

// New loads or initialises a Freezer backed by the given file path.
func New(path string) (*Freezer, error) {
	if path == "" {
		return nil, fmt.Errorf("envfreeze: path must not be empty")
	}
	f := &Freezer{path: path, set: &FrozenSet{}}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return f, nil
	}
	if err != nil {
		return nil, fmt.Errorf("envfreeze: read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, f.set); err != nil {
		return nil, fmt.Errorf("envfreeze: parse %s: %w", path, err)
	}
	return f, nil
}

// Freeze adds keys to the frozen set and persists the result.
func (f *Freezer) Freeze(keys []string, comment string, clock func() time.Time) error {
	if clock == nil {
		clock = time.Now
	}
	existing := toSet(f.set.Keys)
	for _, k := range keys {
		if k == "" {
			continue
		}
		existing[k] = struct{}{}
	}
	merged := make([]string, 0, len(existing))
	for k := range existing {
		merged = append(merged, k)
	}
	sort.Strings(merged)
	f.set = &FrozenSet{Keys: merged, FrozenAt: clock(), Comment: comment}
	return f.save()
}

// IsFrozen reports whether the given key is currently frozen.
func (f *Freezer) IsFrozen(key string) bool {
	for _, k := range f.set.Keys {
		if k == key {
			return true
		}
	}
	return false
}

// Keys returns all currently frozen keys.
func (f *Freezer) Keys() []string {
	out := make([]string, len(f.set.Keys))
	copy(out, f.set.Keys)
	return out
}

// Unfreeze removes the given key from the frozen set and persists.
func (f *Freezer) Unfreeze(key string) error {
	filtered := f.set.Keys[:0:0]
	for _, k := range f.set.Keys {
		if k != key {
			filtered = append(filtered, k)
		}
	}
	f.set.Keys = filtered
	return f.save()
}

func (f *Freezer) save() error {
	if err := os.MkdirAll(filepath.Dir(f.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(f.set, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.path, data, 0o644)
}

func toSet(keys []string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}
