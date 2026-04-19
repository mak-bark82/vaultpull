package envbackup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write env file: %v", err)
	}
	return p
}

func TestBackup_CreatesBackupFile(t *testing.T) {
	dir := t.TempDir()
	envPath := writeEnvFile(t, dir, ".env", "KEY=value\nFOO=bar\n")

	backupPath, err := Backup(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backupPath == "" {
		t.Fatal("expected non-empty backup path")
	}

	data, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if string(data) != "KEY=value\nFOO=bar\n" {
		t.Errorf("backup content mismatch: %q", string(data))
	}
	if !strings.Contains(filepath.Base(backupPath), ".env.") {
		t.Errorf("backup filename unexpected: %s", backupPath)
	}
	if !strings.HasSuffix(backupPath, ".bak") {
		t.Errorf("expected .bak suffix, got: %s", backupPath)
	}
}

func TestBackup_NonExistentFile_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	backupPath, err := Backup(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backupPath != "" {
		t.Errorf("expected empty path for missing file, got: %s", backupPath)
	}
}

func TestBackup_PreservesPermissions(t *testing.T) {
	dir := t.TempDir()
	envPath := writeEnvFile(t, dir, ".env", "SECRET=abc\n")

	backupPath, err := Backup(envPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(backupPath)
	if err != nil {
		t.Fatalf("stat backup: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}
