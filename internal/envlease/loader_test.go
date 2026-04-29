package envlease_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpull/internal/envlease"
)

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "lease.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeConfigFile: %v", err)
	}
	return p
}

func TestLoadConfig_EmptyPath_ReturnsDefaults(t *testing.T) {
	cfg, err := envlease.LoadConfig("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultTTL != time.Hour {
		t.Errorf("DefaultTTL = %v, want 1h", cfg.DefaultTTL)
	}
}

func TestLoadConfig_Valid(t *testing.T) {
	p := writeConfigFile(t, "path: /tmp/leases.json\ndefault_ttl: 30m\n")
	cfg, err := envlease.LoadConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Path != "/tmp/leases.json" {
		t.Errorf("Path = %q, want /tmp/leases.json", cfg.Path)
	}
	if cfg.DefaultTTL != 30*time.Minute {
		t.Errorf("DefaultTTL = %v, want 30m", cfg.DefaultTTL)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := envlease.LoadConfig("/nonexistent/lease.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfig_ZeroTTL_UsesDefault(t *testing.T) {
	p := writeConfigFile(t, "path: /tmp/leases.json\n")
	cfg, err := envlease.LoadConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultTTL != time.Hour {
		t.Errorf("DefaultTTL = %v, want 1h (default)", cfg.DefaultTTL)
	}
}
