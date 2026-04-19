package config

import (
	"errors"
	"os"
)

// Config holds runtime configuration for vaultpull.
type Config struct {
	VaultAddr  string
	VaultToken string
	AuditLog   string
	MappingFile string
	DryRun     bool
	Overwrite  bool
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	addr := os.Getenv("VAULT_ADDR")
	if addr == "" {
		return nil, errors.New("VAULT_ADDR is required")
	}
	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, errors.New("VAULT_TOKEN is required")
	}
	mappingFile := os.Getenv("VAULTPULL_MAPPING")
	if mappingFile == "" {
		mappingFile = "vaultpull.yaml"
	}
	return &Config{
		VaultAddr:   addr,
		VaultToken:  token,
		AuditLog:    os.Getenv("VAULTPULL_AUDIT_LOG"),
		MappingFile: mappingFile,
		DryRun:      os.Getenv("VAULTPULL_DRY_RUN") == "true",
		Overwrite:   os.Getenv("VAULTPULL_OVERWRITE") == "true",
	}, nil
}
