package envrollback

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Snapshot represents a saved state of an env file.
type Snapshot struct {
	File      string
	Timestamp time.Time
	Data      map[string]string
}

// Rollbacker manages rollback snapshots for env files.
type Rollbacker struct {
	snapshotDir string
}

// New creates a Rollbacker that stores snapshots under snapshotDir.
func New(snapshotDir string) (*Rollbacker, error) {
	if snapshotDir == "" {
		return nil, fmt.Errorf("snapshot directory must not be empty")
	}
	if err := os.MkdirAll(snapshotDir, 0700); err != nil {
		return nil, fmt.Errorf("create snapshot dir: %w", err)
	}
	return &Rollbacker{snapshotDir: snapshotDir}, nil
}

// Save persists a snapshot of the given env map for the target file.
func (r *Rollbacker) Save(targetFile string, data map[string]string, now time.Time) error {
	name := snapshotName(targetFile, now)
	path := filepath.Join(r.snapshotDir, name)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open snapshot file: %w", err)
	}
	defer f.Close()
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if _, err := fmt.Fprintf(f, "%s=%s\n", k, data[k]); err != nil {
			return fmt.Errorf("write snapshot: %w", err)
		}
	}
	return nil
}

// Latest returns the most recent snapshot for targetFile, or nil if none exist.
func (r *Rollbacker) Latest(targetFile string) (*Snapshot, error) {
	entries, err := os.ReadDir(r.snapshotDir)
	if err != nil {
		return nil, fmt.Errorf("read snapshot dir: %w", err)
	}
	prefix := snapshotPrefix(targetFile)
	var best os.DirEntry
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), prefix) {
			if best == nil || e.Name() > best.Name() {
				best = e
			}
		}
	}
	if best == nil {
		return nil, nil
	}
	return r.load(targetFile, filepath.Join(r.snapshotDir, best.Name()))
}

func (r *Rollbacker) load(targetFile, path string) (*Snapshot, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read snapshot: %w", err)
	}
	data := map[string]string{}
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			data[parts[0]] = parts[1]
		}
	}
	name := filepath.Base(path)
	ts := parseTimestamp(name)
	return &Snapshot{File: targetFile, Timestamp: ts, Data: data}, nil
}

func snapshotPrefix(targetFile string) string {
	return strings.ReplaceAll(filepath.Base(targetFile), ".", "_") + "_"
}

func snapshotName(targetFile string, t time.Time) string {
	return fmt.Sprintf("%s%s.snap", snapshotPrefix(targetFile), t.UTC().Format("20060102T150405Z"))
}

func parseTimestamp(name string) time.Time {
	parts := strings.Split(name, "_")
	if len(parts) < 2 {
		return time.Time{}
	}
	raw := strings.TrimSuffix(parts[len(parts)-1], ".snap")
	t, _ := time.Parse("20060102T150405Z", raw)
	return t
}
