package envaudit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envaudit"
)

func TestExportJSON_ValidOutput(t *testing.T) {
	r := envaudit.New(fixedClock())
	r.Record(envaudit.EventAdded, "API_KEY", ".env", "new secret")
	r.Record(envaudit.EventChanged, "DB_PASS", ".env.prod", "")

	var buf bytes.Buffer
	if err := r.ExportJSON(&buf); err != nil {
		t.Fatalf("ExportJSON error: %v", err)
	}

	var entries []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if entries[0]["key"] != "API_KEY" {
		t.Errorf("unexpected key: %v", entries[0]["key"])
	}
	if entries[0]["note"] != "new secret" {
		t.Errorf("expected note to be preserved")
	}
	if _, ok := entries[1]["note"]; ok {
		t.Error("empty note should be omitted from JSON")
	}
}

func TestExportJSON_EmptyRecorder(t *testing.T) {
	r := envaudit.New(nil)
	var buf bytes.Buffer
	if err := r.ExportJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(buf.String()) != "[]" {
		t.Errorf("expected empty JSON array, got: %s", buf.String())
	}
}

func TestExportCSV_ContainsFields(t *testing.T) {
	r := envaudit.New(fixedClock())
	r.Record(envaudit.EventSkipped, "TOKEN", ".env", "no-overwrite")

	var buf bytes.Buffer
	if err := r.ExportCSV(&buf); err != nil {
		t.Fatalf("ExportCSV error: %v", err)
	}
	line := buf.String()
	for _, want := range []string{"skipped", "TOKEN", ".env", "no-overwrite"} {
		if !strings.Contains(line, want) {
			t.Errorf("CSV missing field %q in: %s", want, line)
		}
	}
}

func TestExportCSV_MultipleRows(t *testing.T) {
	r := envaudit.New(fixedClock())
	r.Record(envaudit.EventAdded, "X", ".env", "")
	r.Record(envaudit.EventAdded, "Y", ".env", "")

	var buf bytes.Buffer
	_ = r.ExportCSV(&buf)
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 CSV rows, got %d", len(lines))
	}
}
