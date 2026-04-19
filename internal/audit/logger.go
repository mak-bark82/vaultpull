package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	VaultPath string    `json:"vault_path"`
	EnvFile   string    `json:"env_file"`
	Keys      []string  `json:"keys"`
	Status    string    `json:"status"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes audit entries to a file.
type Logger struct {
	path string
	file *os.File
}

// NewLogger opens or creates the audit log file at the given path.
func NewLogger(path string) (*Logger, error) {
	if path == "" {
		return &Logger{}, nil
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &Logger{path: path, file: f}, nil
}

// Log writes an audit entry. If no file is configured it is a no-op.
func (l *Logger) Log(e Entry) error {
	if l.file == nil {
		return nil
	}
	e.Timestamp = time.Now().UTC()
	line, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.file, "%s\n", line)
	return err
}

// Close closes the underlying log file.
func (l *Logger) Close() error {
	if l.file == nil {
		return nil
	}
	return l.file.Close()
}
