package envttl

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Entry holds a secret key with its expiry metadata.
type Entry struct {
	Key       string    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Store manages TTL records for env secret keys.
type Store struct {
	path    string
	entries map[string]Entry
	clock   func() time.Time
}

// New creates a Store backed by the given file path.
func New(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("envttl: path must not be empty")
	}
	s := &Store{
		path:    path,
		entries: make(map[string]Entry),
		clock:   time.Now,
	}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

// Set registers or updates the TTL for a key.
func (s *Store) Set(key string, ttl time.Duration) error {
	if key == "" {
		return errors.New("envttl: key must not be empty")
	}
	s.entries[key] = Entry{
		Key:       key,
		ExpiresAt: s.clock().Add(ttl),
	}
	return s.save()
}

// IsExpired reports whether the key's TTL has elapsed.
// Keys with no TTL record are considered not expired.
func (s *Store) IsExpired(key string) bool {
	e, ok := s.entries[key]
	if !ok {
		return false
	}
	return s.clock().After(e.ExpiresAt)
}

// Remove deletes the TTL record for a key.
func (s *Store) Remove(key string) error {
	delete(s.entries, key)
	return s.save()
}

// Expired returns all keys whose TTL has elapsed.
func (s *Store) Expired() []string {
	var out []string
	for k := range s.entries {
		if s.IsExpired(k) {
			out = append(out, k)
		}
	}
	return out
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.entries)
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}
