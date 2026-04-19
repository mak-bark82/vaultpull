package envlock

import (
	"encoding/json"
	"os"
	"time"
)

// LockEntry records the last known state of a secret path sync.
type LockEntry struct {
	VaultPath string            `json:"vault_path"`
	EnvFile   string            `json:"env_file"`
	SyncedAt  time.Time         `json:"synced_at"`
	Keys      map[string]string `json:"keys"` // key -> sha256 hash of value
}

// LockFile represents the full lock file structure.
type LockFile struct {
	Entries []LockEntry `json:"entries"`
}

// Load reads a lock file from disk. Returns an empty LockFile if not found.
func Load(path string) (*LockFile, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &LockFile{}, nil
	}
	if err != nil {
		return nil, err
	}
	var lf LockFile
	if err := json.Unmarshal(data, &lf); err != nil {
		return nil, err
	}
	return &lf, nil
}

// Save writes the lock file to disk.
func Save(path string, lf *LockFile) error {
	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// Upsert adds or updates a LockEntry for the given vault path.
func (lf *LockFile) Upsert(entry LockEntry) {
	for i, e := range lf.Entries {
		if e.VaultPath == entry.VaultPath && e.EnvFile == entry.EnvFile {
			lf.Entries[i] = entry
			return
		}
	}
	lf.Entries = append(lf.Entries, entry)
}

// Find returns the LockEntry for the given vault path and env file, if any.
func (lf *LockFile) Find(vaultPath, envFile string) (LockEntry, bool) {
	for _, e := range lf.Entries {
		if e.VaultPath == vaultPath && e.EnvFile == envFile {
			return e, true
		}
	}
	return LockEntry{}, false
}
