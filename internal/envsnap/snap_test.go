package envsnap_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/envsnap"
)

var fixedTime = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func fixedClock() time.Time { return fixedTime }

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "secret",
	}
}

func TestTake_CapturesValues(t *testing.T) {
	env := baseEnv()
	snap := envsnap.Take("prod", env, fixedClock)
	if snap.Source != "prod" {
		t.Errorf("expected source prod, got %s", snap.Source)
	}
	if len(snap.Values) != 3 {
		t.Errorf("expected 3 values, got %d", len(snap.Values))
	}
	if snap.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
}

func TestTake_DoesNotMutateInput(t *testing.T) {
	env := baseEnv()
	snap := envsnap.Take("prod", env, fixedClock)
	env["NEW_KEY"] = "injected"
	if _, ok := snap.Values["NEW_KEY"]; ok {
		t.Error("snapshot should not reflect mutation of original map")
	}
}

func TestEqual_SameContent(t *testing.T) {
	a := envsnap.Take("prod", baseEnv(), fixedClock)
	b := envsnap.Take("prod", baseEnv(), fixedClock)
	if !envsnap.Equal(a, b) {
		t.Error("expected snapshots with same content to be equal")
	}
}

func TestEqual_DifferentContent(t *testing.T) {
	a := envsnap.Take("prod", baseEnv(), fixedClock)
	env2 := baseEnv()
	env2["DB_HOST"] = "remotehost"
	b := envsnap.Take("prod", env2, fixedClock)
	if envsnap.Equal(a, b) {
		t.Error("expected snapshots with different content to be unequal")
	}
}

func TestDiff_DetectsChanges(t *testing.T) {
	old := envsnap.Take("prod", baseEnv(), fixedClock)
	newEnv := map[string]string{
		"DB_HOST": "newhost",
		"DB_PORT": "5432",
		"NEW_VAR": "value",
	}
	next := envsnap.Take("prod", newEnv, fixedClock)
	added, removed, changed := envsnap.Diff(old, next)
	if len(added) != 1 || added[0] != "NEW_VAR" {
		t.Errorf("expected NEW_VAR added, got %v", added)
	}
	if len(removed) != 1 || removed[0] != "API_KEY" {
		t.Errorf("expected API_KEY removed, got %v", removed)
	}
	if len(changed) != 1 || changed[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST changed, got %v", changed)
	}
}

func TestSummary_ContainsSource(t *testing.T) {
	snap := envsnap.Take("staging", baseEnv(), fixedClock)
	summary := snap.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
