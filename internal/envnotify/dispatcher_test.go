package envnotify

import (
	"testing"

	"github.com/dstotijn/vaultpull/internal/diff"
)

func TestDispatch_AddedEmitsInfo(t *testing.T) {
	n := New(nil)
	d := NewDispatcher(n)

	d.Dispatch([]diff.Change{
		{Key: "NEW_KEY", Type: diff.Added, New: "val"},
	})

	events := n.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Level != LevelInfo {
		t.Errorf("expected INFO for added key, got %s", events[0].Level)
	}
	if events[0].Key != "NEW_KEY" {
		t.Errorf("expected key NEW_KEY, got %s", events[0].Key)
	}
}

func TestDispatch_RemovedEmitsWarn(t *testing.T) {
	n := New(nil)
	d := NewDispatcher(n)

	d.Dispatch([]diff.Change{
		{Key: "OLD_KEY", Type: diff.Removed, Old: "v"},
	})

	events := n.Events()
	if events[0].Level != LevelWarn {
		t.Errorf("expected WARN for removed key, got %s", events[0].Level)
	}
}

func TestDispatch_ChangedEmitsWarn(t *testing.T) {
	n := New(nil)
	d := NewDispatcher(n)

	d.Dispatch([]diff.Change{
		{Key: "HOST", Type: diff.Changed, Old: "a", New: "b"},
	})

	events := n.Events()
	if events[0].Level != LevelWarn {
		t.Errorf("expected WARN for changed key, got %s", events[0].Level)
	}
}

func TestDispatch_EmptyChanges(t *testing.T) {
	n := New(nil)
	d := NewDispatcher(n)

	d.Dispatch([]diff.Change{})

	if len(n.Events()) != 0 {
		t.Error("expected no events for empty change set")
	}
}

func TestDispatch_MultipleChanges(t *testing.T) {
	n := New(nil)
	d := NewDispatcher(n)

	d.Dispatch([]diff.Change{
		{Key: "A", Type: diff.Added},
		{Key: "B", Type: diff.Removed},
		{Key: "C", Type: diff.Changed},
	})

	if len(n.Events()) != 3 {
		t.Errorf("expected 3 events, got %d", len(n.Events()))
	}
}
