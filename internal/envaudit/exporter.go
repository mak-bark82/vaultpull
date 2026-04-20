package envaudit

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// JSONEntry is the serialisable form of an Event.
type JSONEntry struct {
	Timestamp string `json:"timestamp"`
	Kind      string `json:"kind"`
	Key       string `json:"key"`
	EnvFile   string `json:"env_file"`
	Note      string `json:"note,omitempty"`
}

// ExportJSON writes all recorded events as a JSON array to w.
func (r *Recorder) ExportJSON(w io.Writer) error {
	events := r.Events()
	entries := make([]JSONEntry, len(events))
	for i, e := range events {
		entries[i] = JSONEntry{
			Timestamp: e.Timestamp.Format(time.RFC3339),
			Kind:      string(e.Kind),
			Key:       e.Key,
			EnvFile:   e.EnvFile,
			Note:      e.Note,
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entries); err != nil {
		return fmt.Errorf("envaudit: export json: %w", err)
	}
	return nil
}

// ExportCSV writes all recorded events as CSV lines (no header) to w.
func (r *Recorder) ExportCSV(w io.Writer) error {
	for _, e := range r.Events() {
		line := fmt.Sprintf("%s,%s,%s,%s,%s\n",
			e.Timestamp.Format(time.RFC3339),
			e.Kind,
			e.Key,
			e.EnvFile,
			e.Note,
		)
		if _, err := fmt.Fprint(w, line); err != nil {
			return fmt.Errorf("envaudit: export csv: %w", err)
		}
	}
	return nil
}
