package envsnap

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const snapExt = ".snap.json"

// Store persists snapshots to a directory on disk.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir, creating it if necessary.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

// Save writes a snapshot to disk using its timestamp as the filename.
func (s *Store) Save(snap *Snapshot) error {
	name := snap.Timestamp.UTC().Format("20060102T150405Z") + "_" + sanitizeName(snap.Source) + snapExt
	path := filepath.Join(s.dir, name)
	data, err := json.Marshal(snap)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// Latest returns the most recently saved snapshot, or nil if none exist.
func (s *Store) Latest() (*Snapshot, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), snapExt) {
			names = append(names, e.Name())
		}
	}
	if len(names) == 0 {
		return nil, nil
	}
	sort.Strings(names)
	return s.load(filepath.Join(s.dir, names[len(names)-1]))
}

// All returns all stored snapshots in chronological order.
func (s *Store) All() ([]*Snapshot, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	var snaps []*Snapshot
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), snapExt) {
			continue
		}
		snap, err := s.load(filepath.Join(s.dir, e.Name()))
		if err != nil {
			return nil, err
		}
		snaps = append(snaps, snap)
	}
	return snaps, nil
}

func (s *Store) load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

func sanitizeName(s string) string {
	r := strings.NewReplacer("/", "_", "\\", "_", " ", "_")
	return r.Replace(s)
}
