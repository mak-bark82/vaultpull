package envquota

import (
	"os"
	"path/filepath"
	"testing"
)

func writeQuotaFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "quota.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write quota file: %v", err)
	}
	return p
}

func TestLoadRule_Valid(t *testing.T) {
	p := writeQuotaFile(t, "max_keys: 20\nmax_key_length: 64\nmax_val_length: 256\n")
	rule, err := LoadRule(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.MaxKeys != 20 {
		t.Errorf("expected MaxKeys=20, got %d", rule.MaxKeys)
	}
	if rule.MaxKeyLength != 64 {
		t.Errorf("expected MaxKeyLength=64, got %d", rule.MaxKeyLength)
	}
	if rule.MaxValLength != 256 {
		t.Errorf("expected MaxValLength=256, got %d", rule.MaxValLength)
	}
}

func TestLoadRule_EmptyPath(t *testing.T) {
	rule, err := LoadRule("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.MaxKeys != 0 || rule.MaxKeyLength != 0 || rule.MaxValLength != 0 {
		t.Error("expected zero-value rule for empty path")
	}
}

func TestLoadRule_MissingFile(t *testing.T) {
	_, err := LoadRule("/nonexistent/quota.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadRule_NegativeValue(t *testing.T) {
	p := writeQuotaFile(t, "max_keys: -1\n")
	_, err := LoadRule(p)
	if err == nil {
		t.Fatal("expected error for negative limit")
	}
}

func TestLoadRule_InvalidYAML(t *testing.T) {
	p := writeQuotaFile(t, ": :: invalid yaml :::")
	_, err := LoadRule(p)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
