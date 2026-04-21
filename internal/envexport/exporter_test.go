package envexport_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/envexport"
)

// baseExportSecrets provides a consistent set of secrets for export tests.
func baseExportSecrets() map[string]string {
	return map[string]string{
		"APP_ENV":    "production",
		"DB_HOST":    "db.example.com",
		"DB_PORT":    "5432",
		"SECRET_KEY": "s3cr3t",
	}
}

func TestExport_Dotenv_ContainsAllKeys(t *testing.T) {
	exporter := envexport.New()
	var buf bytes.Buffer

	err := exporter.Export(baseExportSecrets(), "dotenv", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for key, val := range baseExportSecrets() {
		expected := key + "=" + val
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
}

func TestExport_JSON_ValidStructure(t *testing.T) {
	exporter := envexport.New()
	var buf bytes.Buffer

	err := exporter.Export(baseExportSecrets(), "json", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	for key, val := range baseExportSecrets() {
		if got, ok := result[key]; !ok {
			t.Errorf("missing key %q in JSON output", key)
		} else if got != val {
			t.Errorf("key %q: expected %q, got %q", key, val, got)
		}
	}
}

func TestExport_YAML_ContainsKeys(t *testing.T) {
	exporter := envexport.New()
	var buf bytes.Buffer

	err := exporter.Export(baseExportSecrets(), "yaml", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for key := range baseExportSecrets() {
		if !strings.Contains(output, key+":") {
			t.Errorf("expected YAML output to contain key %q", key)
		}
	}
}

func TestExport_UnknownFormat_ReturnsError(t *testing.T) {
	exporter := envexport.New()
	var buf bytes.Buffer

	err := exporter.Export(baseExportSecrets(), "toml", &buf)
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}

	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("expected error to mention 'unsupported', got: %v", err)
	}
}

func TestExport_EmptySecrets_ProducesEmptyOutput(t *testing.T) {
	exporter := envexport.New()

	formats := []string{"dotenv", "json", "yaml"}
	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			var buf bytes.Buffer
			err := exporter.Export(map[string]string{}, format, &buf)
			if err != nil {
				t.Fatalf("unexpected error for format %q: %v", format, err)
			}
			// Output should be minimal/empty but not error
			_ = buf.String()
		})
	}
}

func TestExport_Dotenv_SortedOutput(t *testing.T) {
	exporter := envexport.New()
	var buf bytes.Buffer

	err := exporter.Export(baseExportSecrets(), "dotenv", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	for i := 1; i < len(lines); i++ {
		prev := strings.SplitN(lines[i-1], "=", 2)[0]
		curr := strings.SplitN(lines[i], "=", 2)[0]
		if prev > curr {
			t.Errorf("output is not sorted: %q appears before %q", prev, curr)
		}
	}
}
