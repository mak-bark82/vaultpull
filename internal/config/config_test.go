package config

import (
	"os"
	"testing"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	t.Setenv(key, value)
}

func setValidEnv(t *testing.T) {
	t.Helper()
	setEnv(t, "VAULT_ADDR", "http://127.0.0.1:8200")
	setEnv(t, "VAULT_TOKEN", "test-token")
	setEnv(t, "VAULTPULL_SECRET_PATH", "secret/data/myapp")
}

func TestLoad_Valid(t *testing.T) {
	setValidEnv(t)

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.VaultAddress != "http://127.0.0.1:8200" {
		t.Errorf("unexpected VaultAddress: %s", cfg.VaultAddress)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default OutputFile '.env', got %s", cfg.OutputFile)
	}
	if cfg.Overwrite != false {
		t.Errorf("expected Overwrite false by default")
	}
}

func TestLoad_MissingVaultAddr(t *testing.T) {
	os.Unsetenv("VAULT_ADDR")
	setEnv(t, "VAULT_TOKEN", "test-token")
	setEnv(t, "VAULTPULL_SECRET_PATH", "secret/data/myapp")

	_, err := Load("")
	if err == nil || err.Error() != "VAULT_ADDR is required" {
		t.Errorf("expected VAULT_ADDR error, got %v", err)
	}
}

func TestLoad_MissingToken(t *testing.T) {
	setEnv(t, "VAULT_ADDR", "http://127.0.0.1:8200")
	os.Unsetenv("VAULT_TOKEN")
	setEnv(t, "VAULTPULL_SECRET_PATH", "secret/data/myapp")

	_, err := Load("")
	if err == nil || err.Error() != "VAULT_TOKEN is required" {
		t.Errorf("expected VAULT_TOKEN error, got %v", err)
	}
}

func TestLoad_OverwriteTrue(t *testing.T) {
	setValidEnv(t)
	setEnv(t, "VAULTPULL_OVERWRITE", "true")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Overwrite {
		t.Errorf("expected Overwrite to be true")
	}
}

func TestLoad_CustomOutputFile(t *testing.T) {
	setValidEnv(t)
	setEnv(t, "VAULTPULL_OUTPUT_FILE", ".env.local")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFile != ".env.local" {
		t.Errorf("expected .env.local, got %s", cfg.OutputFile)
	}
}
