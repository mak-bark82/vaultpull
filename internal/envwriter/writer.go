package envwriter

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Options controls how the .env file is written.
type Options struct {
	// Overwrite existing keys if true; skip them if false.
	Overwrite bool
}

// Write merges secrets into the .env file at filePath.
// It preserves existing entries and appends or updates based on Options.
func Write(filePath string, secrets map[string]string, opts Options) error {
	existing := map[string]string{}

	data, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	var lines []string
	if len(data) > 0 {
		lines = strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	}

	for _, line := range lines {
		if idx := strings.IndexByte(line, '='); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := line[idx+1:]
			existing[key] = val
		}
	}

	for k, v := range secrets {
		if _, found := existing[k]; found && !opts.Overwrite {
			continue
		}
		existing[k] = v
	}

	keys := make([]string, 0, len(existing))
	for k := range existing {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(existing[k])
		sb.WriteByte('\n')
	}

	if err := os.WriteFile(filePath, []byte(sb.String()), 0600); err != nil {
		return fmt.Errorf("failed to write %s: %w", filePath, err)
	}
	return nil
}
