package envtemplate

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ParseTemplate reads a .env.template file and returns a list of required keys.
// Lines starting with '#' are treated as comments and skipped.
// Lines of the form KEY or KEY=default are parsed.
type Entry struct {
	Key     string
	Default string
	HasDefault bool
}

func ParseTemplate(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envtemplate: open %q: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		if key == "" {
			continue
		}
		entry := Entry{Key: key}
		if len(parts) == 2 {
			entry.Default = strings.TrimSpace(parts[1])
			entry.HasDefault = true
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envtemplate: scan %q: %w", path, err)
	}
	return entries, nil
}

// ApplyDefaults fills in default values for keys missing from secrets.
func ApplyDefaults(secrets map[string]string, entries []Entry) map[string]string {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}
	for _, e := range entries {
		if _, ok := result[e.Key]; !ok && e.HasDefault {
			result[e.Key] = e.Default
		}
	}
	return result
}
