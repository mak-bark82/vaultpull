package envSign

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Record holds a stored signature and metadata.
type Record struct {
	File      string    `json:"file"`
	Signature string    `json:"signature"`
	SignedAt  time.Time `json:"signed_at"`
}

// SaveRecord writes a signature record to the given path as JSON.
func SaveRecord(path string, rec Record) error {
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadRecord reads a signature record from the given path.
// Returns an error if the file does not exist or is malformed.
func LoadRecord(path string) (*Record, error) {
	if path == "" {
		return nil, errors.New("envSign: record path must not be empty")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}
