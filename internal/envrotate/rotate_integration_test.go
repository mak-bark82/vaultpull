package envrotate_test

import (
	"testing"
	"time"

	"github.com/nicholasgasior/vaultpull/internal/envrotate"
)

// TestApply_RotatedAtTimestamp verifies rotation records carry a recent timestamp.
func TestApply_RotatedAtTimestamp(t *testing.T) {
	before := time.Now().UTC().Add(-time.Second)
	r := envrotate.New(map[string]string{"KEY": "old"})
	records, _ := r.Apply(map[string]string{"KEY": "new"})
	if len(records) != 1 {
		t.Fatal("expected 1 record")
	}
	if records[0].RotatedAt.Before(before) {
		t.Errorf("RotatedAt %v is before test start %v", records[0].RotatedAt, before)
	}
}

// TestApply_DoesNotMutateInput ensures the original map passed to New is unchanged.
func TestApply_DoesNotMutateInput(t *testing.T) {
	orig := map[string]string{"SECRET": "original"}
	r := envrotate.New(orig)
	_, _ = r.Apply(map[string]string{"SECRET": "mutated"})
	if orig["SECRET"] != "original" {
		t.Errorf("input map was mutated, got %s", orig["SECRET"])
	}
}

// TestApply_EmptyIncoming returns no records and preserves current state.
func TestApply_EmptyIncoming(t *testing.T) {
	current := map[string]string{"A": "1", "B": "2"}
	r := envrotate.New(current)
	records, result := r.Apply(map[string]string{})
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}
	if len(result) != len(current) {
		t.Errorf("result length mismatch: expected %d got %d", len(current), len(result))
	}
}
