package envaudit_test

import (
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/envaudit"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func fixedClock() func() time.Time {
	return func() time.Time { return fixedTime }
}

func TestRecord_StoresEvent(t *testing.T) {
	r := envaudit.New(fixedClock())
	r.Record(envaudit.EventAdded, "DB_HOST", ".env", "")
	events := r.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	e := events[0]
	if e.Kind != envaudit.EventAdded {
		t.Errorf("expected kind added, got %s", e.Kind)
	}
	if e.Key != "DB_HOST" {
		t.Errorf("unexpected key: %s", e.Key)
	}
	if e.EnvFile != ".env" {
		t.Errorf("unexpected env file: %s", e.EnvFile)
	}
}

func TestEvents_ReturnsCopy(t *testing.T) {
	r := envaudit.New(fixedClock())
	r.Record(envaudit.EventChanged, "SECRET", ".env.prod", "rotated")
	events := r.Events()
	events[0].Key = "MUTATED"
	original := r.Events()
	if original[0].Key == "MUTATED" {
		t.Error("Events() should return a copy, not a reference")
	}
}

func TestSummary_NoEvents(t *testing.T) {
	r := envaudit.New(nil)
	got := r.Summary()
	if got != "no audit events recorded" {
		t.Errorf("unexpected summary: %s", got)
	}
}

func TestSummary_ContainsKeyAndKind(t *testing.T) {
	r := envaudit.New(fixedClock())
	r.Record(envaudit.EventRemoved, "OLD_KEY", ".env", "deprecated")
	summary := r.Summary()
	if !strings.Contains(summary, "removed") {
		t.Error("summary should contain kind 'removed'")
	}
	if !strings.Contains(summary, "OLD_KEY") {
		t.Error("summary should contain key name")
	}
	if !strings.Contains(summary, "deprecated") {
		t.Error("summary should contain note")
	}
}

func TestSummary_MultipleEvents(t *testing.T) {
	r := envaudit.New(fixedClock())
	r.Record(envaudit.EventAdded, "A", ".env", "")
	r.Record(envaudit.EventChanged, "B", ".env", "")
	r.Record(envaudit.EventSkipped, "C", ".env", "no-overwrite")
	lines := strings.Split(r.Summary(), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}
