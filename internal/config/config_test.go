package config

import (
	"os"
	"testing"
)

func setEnv(t *testing.T, key, val string) {
	t.Helper()
	t.Setenv(key, val)
}

func setValidEnv(t *testing.T) {
	setEnv(t, "VAULT_ADDR", "http://127.0.0.1:8200")
	setEnv(t, "VAULT_TOKEN", "root")
}

func TestLoad_Valid(t *testing.T) {
	setValidEnv(t)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("unexpected addr: %s", cfg.VaultAddr)
	}
	if cfg.MappingFile != "vaultpull.yaml" {
		t.Errorf("expected default mapping file, got: %s", cfg.MappingFile)
	}
}

func TestLoad_MissingVaultAddr(t *testing.T) {
	os.Unsetenv("VAULT_ADDR")
	setEnv(t, "VAULT_TOKEN", "root")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing VAULT_ADDR")
	}
}

func TestLoad_MissingToken(t *testing.T) {
	setEnv(t, "VAULT_ADDR", "http://127.0.0.1:8200")
	os.Unsetenv("VAULT_TOKEN")
	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing VAULT_TOKEN")
	}
}

func TestLoad_DryRun(t *testing.T) {
	setValidEnv(t)
	setEnv(t, "VAULTPULL_DRY_RUN", "true")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestLoad_CustomMappingFile(t *testing.T) {
	setValidEnv(t)
	setEnv(t, "VAULTPULL_MAPPING", "custom.yaml")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MappingFile != "custom.yaml" {
		t.Errorf("expected custom.yaml, got: %s", cfg.MappingFile)
	}
}
