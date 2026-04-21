// Package envexport provides functionality to export environment variables
// from an in-memory map to various output formats (dotenv, JSON, YAML).
package envexport

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents the output format for exported environment variables.
type Format string

const (
	// FormatDotenv writes KEY=VALUE pairs in dotenv format.
	FormatDotenv Format = "dotenv"
	// FormatJSON writes a JSON object of key/value pairs.
	FormatJSON Format = "json"
	// FormatYAML writes a simple YAML mapping of key/value pairs.
	FormatYAML Format = "yaml"
)

// Exporter writes environment variable maps to an io.Writer in a specified format.
type Exporter struct {
	format Format
}

// New creates a new Exporter for the given format.
// Returns an error if the format is not supported.
func New(format Format) (*Exporter, error) {
	switch format {
	case FormatDotenv, FormatJSON, FormatYAML:
		return &Exporter{format: format}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %q", format)
	}
}

// Export writes the provided secrets map to w using the configured format.
// Keys are sorted alphabetically for deterministic output.
func (e *Exporter) Export(w io.Writer, secrets map[string]string) error {
	switch e.format {
	case FormatDotenv:
		return exportDotenv(w, secrets)
	case FormatJSON:
		return exportJSON(w, secrets)
	case FormatYAML:
		return exportYAML(w, secrets)
	default:
		return fmt.Errorf("unsupported export format: %q", e.format)
	}
}

// exportDotenv writes secrets as KEY=VALUE lines, quoting values that contain spaces.
func exportDotenv(w io.Writer, secrets map[string]string) error {
	keys := sortedKeys(secrets)
	for _, k := range keys {
		v := secrets[k]
		if strings.ContainsAny(v, " \t") {
			v = fmt.Sprintf("%q", v)
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, v); err != nil {
			return fmt.Errorf("writing dotenv entry %q: %w", k, err)
		}
	}
	return nil
}

// exportJSON writes secrets as a pretty-printed JSON object.
func exportJSON(w io.Writer, secrets map[string]string) error {
	// Build an ordered representation using sorted keys.
	ordered := make(map[string]string, len(secrets))
	for k, v := range secrets {
		ordered[k] = v
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(ordered); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}
	return nil
}

// exportYAML writes secrets as a simple YAML mapping (key: value).
// Values containing special characters are single-quoted.
func exportYAML(w io.Writer, secrets map[string]string) error {
	keys := sortedKeys(secrets)
	for _, k := range keys {
		v := secrets[k]
		if needsYAMLQuoting(v) {
			v = fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
		}
		if _, err := fmt.Fprintf(w, "%s: %s\n", k, v); err != nil {
			return fmt.Errorf("writing YAML entry %q: %w", k, err)
		}
	}
	return nil
}

// needsYAMLQuoting reports whether a value should be quoted in YAML output.
func needsYAMLQuoting(v string) bool {
	special := []string{":", "#", "{", "}", "[", "]", ",", "&", "*", "?", "|", "-", "<", ">", "=", "!", "%", "@", "`"}
	for _, ch := range special {
		if strings.Contains(v, ch) {
			return true
		}
	}
	return strings.ContainsAny(v, " \t\n") || v == ""
}

// sortedKeys returns the keys of m in ascending alphabetical order.
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
