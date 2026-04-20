package envwatch_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/envwatch"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
}

func TestWatch_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	writeFile(t, path, "KEY=original\n")

	w := envwatch.New(20 * time.Millisecond)
	if err := w.Watch(path); err != nil {
		t.Fatalf("Watch: %v", err)
	}
	w.Start()
	defer w.Stop()

	time.Sleep(30 * time.Millisecond)
	writeFile(t, path, "KEY=changed\n")

	select {
	case ev := <-w.Events:
		if ev.Path != path {
			t.Errorf("expected path %q, got %q", path, ev.Path)
		}
		if ev.OldHash == ev.NewHash {
			t.Error("expected hashes to differ")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for ChangeEvent")
	}
}

func TestWatch_NoEventWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	writeFile(t, path, "KEY=stable\n")

	w := envwatch.New(20 * time.Millisecond)
	if err := w.Watch(path); err != nil {
		t.Fatalf("Watch: %v", err)
	}
	w.Start()
	defer w.Stop()

	select {
	case ev := <-w.Events:
		t.Errorf("unexpected event: %+v", ev)
	case <-time.After(120 * time.Millisecond):
		// expected: no change
	}
}

func TestWatch_NonExistentFile_ReturnsError(t *testing.T) {
	w := envwatch.New(50 * time.Millisecond)
	err := w.Watch("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestWatch_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	pathA := filepath.Join(dir, "a.env")
	pathB := filepath.Join(dir, "b.env")
	writeFile(t, pathA, "A=1\n")
	writeFile(t, pathB, "B=2\n")

	w := envwatch.New(20 * time.Millisecond)
	for _, p := range []string{pathA, pathB} {
		if err := w.Watch(p); err != nil {
			t.Fatalf("Watch %q: %v", p, err)
		}
	}
	w.Start()
	defer w.Stop()

	time.Sleep(30 * time.Millisecond)
	writeFile(t, pathB, "B=updated\n")

	select {
	case ev := <-w.Events:
		if ev.Path != pathB {
			t.Errorf("expected event for %q, got %q", pathB, ev.Path)
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for ChangeEvent")
	}
}
