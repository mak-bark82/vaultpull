package envlock_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/envlock"
)

func TestLoad_NonExistent(t *testing.T) {
	lf, err := envlock.Load("/tmp/does_not_exist_vaultpull.lock")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(lf.Entries) != 0 {
		t.Errorf("expected empty entries")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vaultpull.lock")

	lf := &envlock.LockFile{}
	lf.Upsert(envlock.LockEntry{
		VaultPath: "secret/app",
		EnvFile:   ".env",
		SyncedAt:  time.Now().UTC().Truncate(time.Second),
		Keys:      map[string]string{"DB_PASS": "abc123hash"},
	})

	if err := envlock.Save(path, lf); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := envlock.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].VaultPath != "secret/app" {
		t.Errorf("unexpected vault path: %s", loaded.Entries[0].VaultPath)
	}
}

func TestUpsert_UpdatesExisting(t *testing.T) {
	lf := &envlock.LockFile{}
	lf.Upsert(envlock.LockEntry{VaultPath: "secret/app", EnvFile: ".env", Keys: map[string]string{"A": "1"}})
	lf.Upsert(envlock.LockEntry{VaultPath: "secret/app", EnvFile: ".env", Keys: map[string]string{"A": "2"}})
	if len(lf.Entries) != 1 {
		t.Errorf("expected 1 entry after upsert, got %d", len(lf.Entries))
	}
	if lf.Entries[0].Keys["A"] != "2" {
		t.Errorf("expected updated value")
	}
}

func TestFind(t *testing.T) {
	lf := &envlock.LockFile{}
	lf.Upsert(envlock.LockEntry{VaultPath: "secret/app", EnvFile: ".env"})

	e, ok := lf.Find("secret/app", ".env")
	if !ok {
		t.Fatal("expected to find entry")
	}
	if e.VaultPath != "secret/app" {
		t.Errorf("unexpected vault path")
	}

	_, ok = lf.Find("secret/other", ".env")
	if ok {
		t.Error("expected not found")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.lock")
	_ = os.WriteFile(path, []byte("not json"), 0600)
	_, err := envlock.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
