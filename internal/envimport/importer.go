package envimport

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Format represents a supported import file format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
)

// Options controls import behaviour.
type Options struct {
	Format    Format
	Overwrite bool
}

// Import reads key-value pairs from src in the given format and merges them
// into dst, returning the resulting map. If Overwrite is false, existing keys
// in dst are preserved.
func Import(src string, dst map[string]string, opts Options) (map[string]string, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, fmt.Errorf("envimport: open %q: %w", src, err)
	}
	defer f.Close()

	var incoming map[string]string
	switch opts.Format {
	case FormatDotenv, "":
		incoming, err = parseDotenv(f)
	default:
		return nil, fmt.Errorf("envimport: unsupported format %q", opts.Format)
	}
	if err != nil {
		return nil, fmt.Errorf("envimport: parse: %w", err)
	}

	result := make(map[string]string, len(dst))
	for k, v := range dst {
		result[k] = v
	}
	for k, v := range incoming {
		if _, exists := result[k]; exists && !opts.Overwrite {
			continue
		}
		result[k] = v
	}
	return result, nil
}

// stripQuotes removes a matching pair of surrounding single or double quotes
// from s, if present. It does not strip mismatched or nested quotes.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func parseDotenv(f *os.File) (map[string]string, error) {
	out := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := stripQuotes(strings.TrimSpace(parts[1]))
		if key == "" {
			continue
		}
		out[key] = val
	}
	return out, scanner.Err()
}
