package envdrift_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envdrift"
)

var baseVault = map[string]string{
	"DB_HOST": "db.prod",
	"DB_PASS": "s3cr3t",
	"API_KEY": "key123",
}

func TestDetect_NoChanges(t *testing.T) {
	local := map[string]string{
		"DB_HOST": "db.prod",
		"DB_PASS": "s3cr3t",
		"API_KEY": "key123",
	}
	report := envdrift.Detect(baseVault, local)
	if report.HasDrift() {
		t.Errorf("expected no drift, got: %s", report.Summary())
	}
}

func TestDetect_Changed(t *testing.T) {
	local := map[string]string{
		"DB_HOST": "db.prod",
		"DB_PASS": "old-pass",
		"API_KEY": "key123",
	}
	report := envdrift.Detect(baseVault, local)
	if !report.HasDrift() {
		t.Fatal("expected drift")
	}
	found := false
	for _, e := range report.Entries {
		if e.Key == "DB_PASS" && e.Status == envdrift.StatusChanged {
			found = true
		}
	}
	if !found {
		t.Error("expected DB_PASS to be StatusChanged")
	}
}

func TestDetect_Missing(t *testing.T) {
	local := map[string]string{
		"DB_HOST": "db.prod",
		"API_KEY": "key123",
	}
	report := envdrift.Detect(baseVault, local)
	for _, e := range report.Entries {
		if e.Key == "DB_PASS" && e.Status != envdrift.StatusMissing {
			t.Errorf("expected DB_PASS StatusMissing, got %v", e.Status)
		}
	}
}

func TestDetect_Extra(t *testing.T) {
	local := map[string]string{
		"DB_HOST":  "db.prod",
		"DB_PASS":  "s3cr3t",
		"API_KEY":  "key123",
		"EXTRA_VAR": "surprise",
	}
	report := envdrift.Detect(baseVault, local)
	found := false
	for _, e := range report.Entries {
		if e.Key == "EXTRA_VAR" && e.Status == envdrift.StatusExtra {
			found = true
		}
	}
	if !found {
		t.Error("expected EXTRA_VAR to be StatusExtra")
	}
}

func TestSummary_NoDrift(t *testing.T) {
	report := envdrift.Detect(baseVault, baseVault)
	if got := report.Summary(); got != "no drift detected" {
		t.Errorf("unexpected summary: %s", got)
	}
}

func TestSummary_WithDrift(t *testing.T) {
	local := map[string]string{"DB_HOST": "other"}
	report := envdrift.Detect(baseVault, local)
	s := report.Summary()
	if s == "no drift detected" {
		t.Error("expected drift summary")
	}
}
