package envsnap_test

import (
	"os"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/envsnap"
)

func TestNewStore_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	subdir := dir + "/snaps"
	_, err := envsnap.NewStore(subdir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(subdir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestSaveAndLatest_Roundtrip(t *testing.T) {
	store, _ := envsnap.NewStore(t.TempDir())
	snap := envsnap.Take("prod", baseEnv(), fixedClock)
	if err := store.Save(snap); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	loaded, err := store.Latest()
	if err != nil {
		t.Fatalf("latest failed: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if loaded.Source != snap.Source {
		t.Errorf("source mismatch: got %s", loaded.Source)
	}
	if loaded.Checksum != snap.Checksum {
		t.Errorf("checksum mismatch")
	}
}

func TestLatest_NilWhenEmpty(t *testing.T) {
	store, _ := envsnap.NewStore(t.TempDir())
	snap, err := store.Latest()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap != nil {
		t.Error("expected nil for empty store")
	}
}

func TestAll_ReturnsSortedSnapshots(t *testing.T) {
	store, _ := envsnap.NewStore(t.TempDir())
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	snap1 := envsnap.Take("prod", baseEnv(), func() time.Time { return t1 })
	snap2 := envsnap.Take("prod", baseEnv(), func() time.Time { return t2 })
	_ = store.Save(snap1)
	_ = store.Save(snap2)
	all, err := store.All()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(all))
	}
	if !all[0].Timestamp.Before(all[1].Timestamp) {
		t.Error("expected snapshots in chronological order")
	}
}

func TestAll_EmptyStore(t *testing.T) {
	store, _ := envsnap.NewStore(t.TempDir())
	all, err := store.All()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 0 {
		t.Errorf("expected 0 snapshots, got %d", len(all))
	}
}
