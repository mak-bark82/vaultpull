package envhash_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envhash"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "secret",
	}
}

func TestHash_IsStable(t *testing.T) {
	h := envhash.New()
	env := baseEnv()

	first := h.Hash(env)
	second := h.Hash(env)

	if first != second {
		t.Errorf("expected stable hash, got %q and %q", first, second)
	}
}

func TestHash_OrderIndependent(t *testing.T) {
	h := envhash.New()

	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAR": "2", "FOO": "1"}

	if h.Hash(a) != h.Hash(b) {
		t.Error("hash should be order-independent")
	}
}

func TestHash_DiffersOnChange(t *testing.T) {
	h := envhash.New()

	original := baseEnv()
	modified := baseEnv()
	modified["API_KEY"] = "rotated"

	if h.Hash(original) == h.Hash(modified) {
		t.Error("expected different hashes after value change")
	}
}

func TestEqual_SameMaps(t *testing.T) {
	h := envhash.New()
	if !h.Equal(baseEnv(), baseEnv()) {
		t.Error("expected equal maps to be equal")
	}
}

func TestEqual_DifferentMaps(t *testing.T) {
	h := envhash.New()
	a := baseEnv()
	b := baseEnv()
	b["NEW_KEY"] = "value"

	if h.Equal(a, b) {
		t.Error("expected unequal maps to not be equal")
	}
}

func TestDiff_DetectsChanged(t *testing.T) {
	h := envhash.New()
	old := map[string]string{"FOO": "old", "BAR": "same"}
	next := map[string]string{"FOO": "new", "BAR": "same"}

	changed := h.Diff(old, next)
	if len(changed) != 1 || changed[0] != "FOO" {
		t.Errorf("expected [FOO], got %v", changed)
	}
}

func TestDiff_DetectsAdded(t *testing.T) {
	h := envhash.New()
	old := map[string]string{"FOO": "1"}
	next := map[string]string{"FOO": "1", "BAR": "2"}

	changed := h.Diff(old, next)
	if len(changed) != 1 || changed[0] != "BAR" {
		t.Errorf("expected [BAR], got %v", changed)
	}
}

func TestDiff_DetectsRemoved(t *testing.T) {
	h := envhash.New()
	old := map[string]string{"FOO": "1", "BAR": "2"}
	next := map[string]string{"FOO": "1"}

	changed := h.Diff(old, next)
	if len(changed) != 1 || changed[0] != "BAR" {
		t.Errorf("expected [BAR], got %v", changed)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	h := envhash.New()
	env := baseEnv()

	changed := h.Diff(env, env)
	if len(changed) != 0 {
		t.Errorf("expected no changes, got %v", changed)
	}
}
