// Package envpin provides functionality to pin secret versions,
// preventing them from being overwritten during future sync operations.
package envpin

import (
	"encoding/json"
	"os"
	"time"
)

// PinEntry records a pinned key and the version it was pinned at.
type PinEntry struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	PinnedAt  time.Time `json:"pinned_at"`
	Comment   string    `json:"comment,omitempty"`
}

// PinFile is the on-disk representation of all pinned keys.
type PinFile struct {
	Pins []PinEntry `json:"pins"`
}

// Pinner manages pinned environment keys.
type Pinner struct {
	file PinFile
}

// Load reads a pin file from the given path. If the path is empty or the file
// does not exist, an empty Pinner is returned without error.
func Load(path string) (*Pinner, error) {
	p := &Pinner{}
	if path == "" {
		return p, nil
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return p, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &p.file); err != nil {
		return nil, err
	}
	return p, nil
}

// Save writes the current pin state to the given path.
func (p *Pinner) Save(path string) error {
	data, err := json.MarshalIndent(p.file, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// Pin adds or updates a pinned entry for the given key.
func (p *Pinner) Pin(key, value, comment string) {
	for i, e := range p.file.Pins {
		if e.Key == key {
			p.file.Pins[i].Value = value
			p.file.Pins[i].PinnedAt = time.Now().UTC()
			p.file.Pins[i].Comment = comment
			return
		}
	}
	p.file.Pins = append(p.file.Pins, PinEntry{
		Key:      key,
		Value:    value,
		PinnedAt: time.Now().UTC(),
		Comment:  comment,
	})
}

// Unpin removes a pinned entry by key. Returns true if the key was found.
func (p *Pinner) Unpin(key string) bool {
	for i, e := range p.file.Pins {
		if e.Key == key {
			p.file.Pins = append(p.file.Pins[:i], p.file.Pins[i+1:]...)
			return true
		}
	}
	return false
}

// IsPinned reports whether the given key is currently pinned.
func (p *Pinner) IsPinned(key string) bool {
	for _, e := range p.file.Pins {
		if e.Key == key {
			return true
		}
	}
	return false
}

// Apply returns a new map based on incoming, but replaces values for pinned
// keys with their pinned values, preserving the pin contract.
func (p *Pinner) Apply(incoming map[string]string) map[string]string {
	out := make(map[string]string, len(incoming))
	for k, v := range incoming {
		out[k] = v
	}
	for _, e := range p.file.Pins {
		out[e.Key] = e.Value
	}
	return out
}

// Entries returns a copy of all current pin entries.
func (p *Pinner) Entries() []PinEntry {
	copy := make([]PinEntry, len(p.file.Pins))
	for i, e := range p.file.Pins {
		copy[i] = e
	}
	return copy
}
