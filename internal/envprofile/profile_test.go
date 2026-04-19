package envprofile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/envprofile"
)

func writeProfileFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "profiles.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writing profile file: %v", err)
	}
	return p
}

func TestLoadProfiles_Valid(t *testing.T) {
	path := writeProfileFile(t, `
profiles:
  dev:
    name: dev
    vault_prefix: secret/dev
    env_file: .env.dev
  prod:
    name: prod
    vault_prefix: secret/prod
    env_file: .env.prod
`)
	ps, err := envprofile.LoadProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ps.Profiles) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(ps.Profiles))
	}
}

func TestLoadProfiles_EmptyPath(t *testing.T) {
	_, err := envprofile.LoadProfiles("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestLoadProfiles_MissingVaultPrefix(t *testing.T) {
	path := writeProfileFile(t, `
profiles:
  dev:
    name: dev
    env_file: .env.dev
`)
	_, err := envprofile.LoadProfiles(path)
	if err == nil {
		t.Fatal("expected error for missing vault_prefix")
	}
}

func TestGet_Found(t *testing.T) {
	path := writeProfileFile(t, `
profiles:
  staging:
    name: staging
    vault_prefix: secret/staging
    env_file: .env.staging
    overrides:
      LOG_LEVEL: debug
`)
	ps, err := envprofile.LoadProfiles(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p, err := ps.Get("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.VaultPrefix != "secret/staging" {
		t.Errorf("expected vault_prefix secret/staging, got %s", p.VaultPrefix)
	}
	if p.Overrides["LOG_LEVEL"] != "debug" {
		t.Errorf("expected override LOG_LEVEL=debug")
	}
}

func TestGet_NotFound(t *testing.T) {
	path := writeProfileFile(t, `
profiles:
  dev:
    name: dev
    vault_prefix: secret/dev
    env_file: .env.dev
`)
	ps, _ := envprofile.LoadProfiles(path)
	_, err := ps.Get("prod")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}
