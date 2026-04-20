package envdiff

import (
	"bytes"
	"strings"
	"testing"
)

func baseOld() map[string]string {
	return map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
		"REMOVED_KEY": "gone",
	}
}

func baseNew() map[string]string {
	return map[string]string{
		"HOST":    "prod.example.com",
		"PORT":    "5432",
		"NEW_KEY": "hello",
	}
}

func TestDiff_DetectsAdded(t *testing.T) {
	r := Diff(baseOld(), baseNew())
	for _, c := range r.Changes {
		if c.Key == "NEW_KEY" && c.Kind == Added {
			return
		}
	}
	t.Error("expected NEW_KEY to be detected as added")
}

func TestDiff_DetectsRemoved(t *testing.T) {
	r := Diff(baseOld(), baseNew())
	for _, c := range r.Changes {
		if c.Key == "REMOVED_KEY" && c.Kind == Removed {
			return
		}
	}
	t.Error("expected REMOVED_KEY to be detected as removed")
}

func TestDiff_DetectsChanged(t *testing.T) {
	r := Diff(baseOld(), baseNew())
	for _, c := range r.Changes {
		if c.Key == "HOST" && c.Kind == Changed {
			if c.OldValue != "localhost" || c.NewValue != "prod.example.com" {
				t.Errorf("unexpected values: old=%q new=%q", c.OldValue, c.NewValue)
			}
			return
		}
	}
	t.Error("expected HOST to be detected as changed")
}

func TestDiff_UnchangedNotReported(t *testing.T) {
	r := Diff(baseOld(), baseNew())
	for _, c := range r.Changes {
		if c.Key == "PORT" {
			t.Errorf("PORT should not appear in changes, got kind=%s", c.Kind)
		}
	}
}

func TestDiff_EmptyMaps(t *testing.T) {
	r := Diff(map[string]string{}, map[string]string{})
	if r.HasChanges() {
		t.Error("expected no changes for empty maps")
	}
}

func TestPrint_ShowsSymbols(t *testing.T) {
	r := Diff(baseOld(), baseNew())
	var buf bytes.Buffer
	r.Print(&buf)
	out := buf.String()

	if !strings.Contains(out, "+ NEW_KEY") {
		t.Error("expected '+' for added key")
	}
	if !strings.Contains(out, "- REMOVED_KEY") {
		t.Error("expected '-' for removed key")
	}
	if !strings.Contains(out, "~ HOST") {
		t.Error("expected '~' for changed key")
	}
}

func TestPrint_NoChanges(t *testing.T) {
	r := Diff(map[string]string{"A": "1"}, map[string]string{"A": "1"})
	var buf bytes.Buffer
	r.Print(&buf)
	if !strings.Contains(buf.String(), "No changes") {
		t.Error("expected 'No changes' message")
	}
}
