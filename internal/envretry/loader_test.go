package envretry_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/vaultpull/internal/envretry"
)

func writePolicyFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "retry.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write policy file: %v", err)
	}
	return p
}

func TestLoadPolicy_EmptyPath(t *testing.T) {
	p, err := envretry.LoadPolicy("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	def := envretry.DefaultPolicy()
	if p.MaxAttempts != def.MaxAttempts {
		t.Errorf("expected MaxAttempts=%d, got %d", def.MaxAttempts, p.MaxAttempts)
	}
}

func TestLoadPolicy_Valid(t *testing.T) {
	path := writePolicyFile(t, "max_attempts: 5\ndelay_ms: 100\nmultiplier: 1.5\n")
	p, err := envretry.LoadPolicy(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.MaxAttempts != 5 {
		t.Errorf("expected 5, got %d", p.MaxAttempts)
	}
	if p.Multiplier != 1.5 {
		t.Errorf("expected 1.5, got %f", p.Multiplier)
	}
}

func TestLoadPolicy_MissingFile(t *testing.T) {
	_, err := envretry.LoadPolicy("/nonexistent/retry.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadPolicy_InvalidMaxAttempts(t *testing.T) {
	path := writePolicyFile(t, "max_attempts: 0\ndelay_ms: 50\nmultiplier: 2.0\n")
	_, err := envretry.LoadPolicy(path)
	if err == nil {
		t.Fatal("expected error for max_attempts=0")
	}
}

func TestLoadPolicy_DefaultMultiplier(t *testing.T) {
	path := writePolicyFile(t, "max_attempts: 2\ndelay_ms: 50\nmultiplier: 0\n")
	p, err := envretry.LoadPolicy(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Multiplier != 1.0 {
		t.Errorf("expected multiplier defaulted to 1.0, got %f", p.Multiplier)
	}
}
