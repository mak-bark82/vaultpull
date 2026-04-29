package envlease_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/envlease"
)

var (
	fixedNow  = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	fixedClock = func() time.Time { return fixedNow }
)

func tempManager(t *testing.T) *envlease.Manager {
	t.Helper()
	path := filepath.Join(t.TempDir(), "leases.json")
	m, err := envlease.New(path, fixedClock)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return m
}

func TestAcquire_GrantsLease(t *testing.T) {
	m := tempManager(t)
	if err := m.Acquire("DB_PASS", "alice", time.Hour); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	l, ok := m.Get("DB_PASS")
	if !ok {
		t.Fatal("expected lease to exist")
	}
	if l.Owner != "alice" {
		t.Errorf("owner = %q, want alice", l.Owner)
	}
}

func TestAcquire_BlocksWhenActive(t *testing.T) {
	m := tempManager(t)
	_ = m.Acquire("DB_PASS", "alice", time.Hour)
	err := m.Acquire("DB_PASS", "bob", time.Hour)
	if err == nil {
		t.Fatal("expected error acquiring already-leased key")
	}
}

func TestAcquire_AllowsAfterExpiry(t *testing.T) {
	now := fixedNow
	clock := func() time.Time { return now }
	path := filepath.Join(t.TempDir(), "leases.json")
	m, _ := envlease.New(path, clock)
	_ = m.Acquire("DB_PASS", "alice", time.Minute)
	// advance past expiry
	now = now.Add(2 * time.Minute)
	if err := m.Acquire("DB_PASS", "bob", time.Hour); err != nil {
		t.Fatalf("expected success after expiry, got: %v", err)
	}
}

func TestRelease_RemovesLease(t *testing.T) {
	m := tempManager(t)
	_ = m.Acquire("API_KEY", "alice", time.Hour)
	if err := m.Release("API_KEY", "alice"); err != nil {
		t.Fatalf("Release: %v", err)
	}
	_, ok := m.Get("API_KEY")
	if ok {
		t.Error("expected lease to be removed")
	}
}

func TestRelease_WrongOwner_ReturnsError(t *testing.T) {
	m := tempManager(t)
	_ = m.Acquire("API_KEY", "alice", time.Hour)
	if err := m.Release("API_KEY", "bob"); err == nil {
		t.Fatal("expected error releasing with wrong owner")
	}
}

func TestPurgeExpired_RemovesStaleLeases(t *testing.T) {
	now := fixedNow
	clock := func() time.Time { return now }
	path := filepath.Join(t.TempDir(), "leases.json")
	m, _ := envlease.New(path, clock)
	_ = m.Acquire("OLD_KEY", "alice", time.Minute)
	_ = m.Acquire("NEW_KEY", "bob", time.Hour)
	now = now.Add(2 * time.Minute)
	if err := m.PurgeExpired(); err != nil {
		t.Fatalf("PurgeExpired: %v", err)
	}
	if _, ok := m.Get("OLD_KEY"); ok {
		t.Error("expected OLD_KEY to be purged")
	}
	if _, ok := m.Get("NEW_KEY"); !ok {
		t.Error("expected NEW_KEY to remain")
	}
}

func TestNew_EmptyPath_InMemory(t *testing.T) {
	m, err := envlease.New("", nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := m.Acquire("X", "alice", time.Hour); err != nil {
		t.Fatalf("Acquire: %v", err)
	}
}
