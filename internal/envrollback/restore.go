package envrollback

import (
	"fmt"
	"os"
	"sort"
)

// RestoreResult summarises what was written during a restore.
type RestoreResult struct {
	File    string
	Written int
}

// Restore writes the given snapshot's data back to targetFile, replacing its
// contents. The caller is responsible for taking a backup before calling this.
func Restore(targetFile string, snap *Snapshot) (*RestoreResult, error) {
	if snap == nil {
		return nil, fmt.Errorf("cannot restore from nil snapshot")
	}
	if targetFile == "" {
		return nil, fmt.Errorf("target file must not be empty")
	}
	f, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("open target file: %w", err)
	}
	defer f.Close()

	keys := make([]string, 0, len(snap.Data))
	for k := range snap.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if _, err := fmt.Fprintf(f, "%s=%s\n", k, snap.Data[k]); err != nil {
			return nil, fmt.Errorf("write key %q: %w", k, err)
		}
	}
	return &RestoreResult{File: targetFile, Written: len(keys)}, nil
}
