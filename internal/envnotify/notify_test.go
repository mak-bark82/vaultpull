package envnotify

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNotify_WritesFormattedLine(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	n.clock = fixedClock(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	n.Notify(LevelInfo, "DB_HOST", "value changed")

	out := buf.String()
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected INFO in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected key in output, got: %s", out)
	}
	if !strings.Contains(out, "value changed") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestNotify_RecordsEvents(t *testing.T) {
	n := New(nil)
	n.Notify(LevelInfo, "A", "added")
	n.Notify(LevelWarn, "B", "changed")

	events := n.Events()
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Level != LevelInfo {
		t.Errorf("expected INFO, got %s", events[0].Level)
	}
	if events[1].Key != "B" {
		t.Errorf("expected key B, got %s", events[1].Key)
	}
}

func TestEvents_ReturnsCopy(t *testing.T) {
	n := New(nil)
	n.Notify(LevelError, "X", "missing")

	e1 := n.Events()
	e1[0].Key = "MUTATED"
	e2 := n.Events()

	if e2[0].Key == "MUTATED" {
		t.Error("Events() should return a copy, not a reference")
	}
}

func TestSummary_NoEvents(t *testing.T) {
	n := New(nil)
	if got := n.Summary(); got != "no events" {
		t.Errorf("expected 'no events', got %q", got)
	}
}

func TestSummary_CountsByLevel(t *testing.T) {
	n := New(nil)
	n.Notify(LevelInfo, "A", "ok")
	n.Notify(LevelInfo, "B", "ok")
	n.Notify(LevelWarn, "C", "warn")
	n.Notify(LevelError, "D", "err")

	s := n.Summary()
	if !strings.Contains(s, "INFO=2") {
		t.Errorf("expected INFO=2 in summary, got: %s", s)
	}
	if !strings.Contains(s, "WARN=1") {
		t.Errorf("expected WARN=1 in summary, got: %s", s)
	}
	if !strings.Contains(s, "ERROR=1") {
		t.Errorf("expected ERROR=1 in summary, got: %s", s)
	}
}
