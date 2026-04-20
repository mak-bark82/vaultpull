package envrollback

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew_EmptyDir(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty snapshot dir")
	}
}

func TestNew_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "snapshots")
	r, err := New(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil rollbacker")
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("snapshot dir was not created")
	}
}

func TestSaveAndLatest_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	r, _ := New(dir)
	data := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	if err := r.Save(".env", data, now); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	snap, err := r.Latest(".env")
	if err != nil {
		t.Fatalf("Latest failed: %v", err)
	}
	if snap == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if snap.Data["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", snap.Data["DB_HOST"])
	}
	if snap.Data["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", snap.Data["PORT"])
	}
}

func TestLatest_ReturnsNilWhenNoSnapshots(t *testing.T) {
	dir := t.TempDir()
	r, _ := New(dir)
	snap, err := r.Latest(".env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap != nil {
		t.Fatalf("expected nil snapshot, got %+v", snap)
	}
}

func TestLatest_ReturnsMostRecent(t *testing.T) {
	dir := t.TempDir()
	r, _ := New(dir)
	old := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	new_ := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	_ = r.Save(".env", map[string]string{"KEY": "old"}, old)
	_ = r.Save(".env", map[string]string{"KEY": "new"}, new_)

	snap, err := r.Latest(".env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Data["KEY"] != "new" {
		t.Errorf("expected KEY=new, got %q", snap.Data["KEY"])
	}
}

func TestSave_PreservesTimestamp(t *testing.T) {
	dir := t.TempDir()
	r, _ := New(dir)
	now := time.Date(2024, 3, 15, 9, 30, 0, 0, time.UTC)
	_ = r.Save(".env", map[string]string{"X": "1"}, now)
	snap, _ := r.Latest(".env")
	if snap.Timestamp != now {
		t.Errorf("expected timestamp %v, got %v", now, snap.Timestamp)
	}
}
