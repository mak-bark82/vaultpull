package dryrun

import (
	"bytes"
	"strings"
	"testing"
)

func TestReport_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	existing := map[string]string{"KEY": "val"}
	incoming := map[string]string{"KEY": "val"}
	r.Report(".env", existing, incoming)
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes', got: %s", buf.String())
	}
}

func TestReport_Added(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	r.Report(".env", map[string]string{}, map[string]string{"NEW_KEY": "value"})
	out := buf.String()
	if !strings.Contains(out, "+ NEW_KEY") {
		t.Errorf("expected added key, got: %s", out)
	}
}

func TestReport_Removed(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	r.Report(".env", map[string]string{"OLD": "v"}, map[string]string{})
	out := buf.String()
	if !strings.Contains(out, "- OLD") {
		t.Errorf("expected removed key, got: %s", out)
	}
}

func TestReport_Changed(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	r.Report(".env", map[string]string{"K": "old"}, map[string]string{"K": "new"})
	out := buf.String()
	if !strings.Contains(out, "~ K") {
		t.Errorf("expected changed key, got: %s", out)
	}
}

func TestReport_NilWriter(t *testing.T) {
	// Should not panic when out is nil (defaults to stdout)
	r := NewReporter(nil)
	if r.out == nil {
		t.Error("expected non-nil writer")
	}
}

func TestReport_Summary(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	r.Report(".env",
		map[string]string{"A": "1", "B": "2"},
		map[string]string{"A": "changed", "C": "3"},
	)
	out := buf.String()
	if !strings.Contains(out, "added") && !strings.Contains(out, "changed") {
		t.Errorf("expected summary line, got: %s", out)
	}
}
