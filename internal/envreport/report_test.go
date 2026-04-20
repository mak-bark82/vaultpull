package envreport_test

import (
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/envreport"
)

func baseEntries() []envreport.Entry {
	return []envreport.Entry{
		{Key: "DB_HOST", Status: "added", Source: "secret/app"},
		{Key: "DB_PASS", Status: "changed", Source: "secret/app"},
		{Key: "OLD_KEY", Status: "removed", Source: "secret/app"},
		{Key: "APP_ENV", Status: "unchanged", Source: "secret/app"},
	}
}

func TestNew_SetsTimestampAndFile(t *testing.T) {
	before := time.Now().UTC()
	r := envreport.New(".env", baseEntries())
	after := time.Now().UTC()

	if r.EnvFile != ".env" {
		t.Errorf("expected env file '.env', got %q", r.EnvFile)
	}
	if r.Timestamp.Before(before) || r.Timestamp.After(after) {
		t.Errorf("timestamp out of expected range: %v", r.Timestamp)
	}
}

func TestSummary_CountsCorrectly(t *testing.T) {
	r := envreport.New(".env", baseEntries())
	s := r.Summary()

	if s["added"] != 1 {
		t.Errorf("expected 1 added, got %d", s["added"])
	}
	if s["changed"] != 1 {
		t.Errorf("expected 1 changed, got %d", s["changed"])
	}
	if s["removed"] != 1 {
		t.Errorf("expected 1 removed, got %d", s["removed"])
	}
	if s["unchanged"] != 1 {
		t.Errorf("expected 1 unchanged, got %d", s["unchanged"])
	}
}

func TestRender_ContainsKeyAndStatus(t *testing.T) {
	r := envreport.New(".env", baseEntries())
	var sb strings.Builder
	if err := r.Render(&sb); err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	out := sb.String()

	for _, want := range []string{"DB_HOST", "added", "changed", "removed", "unchanged", ".env", "vaultpull sync report"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q\nGot:\n%s", want, out)
		}
	}
}

func TestRender_SortsKeys(t *testing.T) {
	entries := []envreport.Entry{
		{Key: "Z_KEY", Status: "added", Source: "secret/app"},
		{Key: "A_KEY", Status: "added", Source: "secret/app"},
		{Key: "M_KEY", Status: "added", Source: "secret/app"},
	}
	r := envreport.New(".env", entries)
	var sb strings.Builder
	_ = r.Render(&sb)
	out := sb.String()

	aIdx := strings.Index(out, "A_KEY")
	mIdx := strings.Index(out, "M_KEY")
	zIdx := strings.Index(out, "Z_KEY")
	if !(aIdx < mIdx && mIdx < zIdx) {
		t.Errorf("keys not sorted alphabetically in output")
	}
}

func TestSummary_EmptyEntries(t *testing.T) {
	r := envreport.New(".env", nil)
	s := r.Summary()
	for k, v := range s {
		if v != 0 {
			t.Errorf("expected 0 for %q, got %d", k, v)
		}
	}
}
