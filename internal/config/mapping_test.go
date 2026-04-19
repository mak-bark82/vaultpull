package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/vaultpull/internal/config"
)

func writeMappingFile(t *testing.T, mappings []config.Mapping) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "mappings.json")
	b, _ := json.Marshal(mappings)
	_ = os.WriteFile(p, b, 0644)
	return p
}

func TestLoadMappings_Valid(t *testing.T) {
	p := writeMappingFile(t, []config.Mapping{
		{VaultPath: "secret/data/app", EnvFile: ".env", Overwrite: true},
	})
	mappings, err := config.LoadMappings(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mappings) != 1 {
		t.Fatalf("expected 1 mapping, got %d", len(mappings))
	}
}

func TestLoadMappings_EmptyPath(t *testing.T) {
	mappings, err := config.LoadMappings("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mappings) != 0 {
		t.Fatal("expected empty mappings")
	}
}

func TestLoadMappings_MissingVaultPath(t *testing.T) {
	p := writeMappingFile(t, []config.Mapping{
		{VaultPath: "", EnvFile: ".env"},
	})
	_, err := config.LoadMappings(p)
	if err == nil {
		t.Fatal("expected error for missing vault_path")
	}
}

func TestLoadMappings_MissingEnvFile(t *testing.T) {
	p := writeMappingFile(t, []config.Mapping{
		{VaultPath: "secret/data/app", EnvFile: ""},
	})
	_, err := config.LoadMappings(p)
	if err == nil {
		t.Fatal("expected error for missing env_file")
	}
}

func TestLoadMappings_FileNotFound(t *testing.T) {
	_, err := config.LoadMappings("/nonexistent/path/mappings.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
