package syncer

import (
	"fmt"

	"github.com/example/vaultpull/internal/audit"
	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/envwriter"
	"github.com/example/vaultpull/internal/vault"
)

// VaultReader abstracts secret reading from Vault.
type VaultReader interface {
	ReadSecrets(path string) (map[string]string, error)
}

// Syncer orchestrates pulling secrets and writing env files.
type Syncer struct {
	client  VaultReader
	logger  *audit.Logger
}

// New creates a Syncer with the provided Vault client and audit logger.
func New(c VaultReader, l *audit.Logger) *Syncer {
	if l == nil {
		l = &audit.Logger{}
	}
	return &Syncer{client: c, logger: l}
}

// Run iterates over mappings, pulls secrets, and writes env files.
func (s *Syncer) Run(mappings []config.Mapping, overwrite bool) error {
	for _, m := range mappings {
		secrets, err := s.client.ReadSecrets(m.VaultPath)
		entry := audit.Entry{
			VaultPath: m.VaultPath,
			EnvFile:   m.EnvFile,
		}
		if err != nil {
			entry.Status = "error"
			entry.Message = err.Error()
			_ = s.logger.Log(entry)
			return fmt.Errorf("syncer: read %q: %w", m.VaultPath, err)
		}
		keys := make([]string, 0, len(secrets))
		for k := range secrets {
			keys = append(keys, k)
		}
		if err := envwriter.Write(m.EnvFile, secrets, overwrite); err != nil {
			entry.Status = "error"
			entry.Message = err.Error()
			_ = s.logger.Log(entry)
			return fmt.Errorf("syncer: write %q: %w", m.EnvFile, err)
		}
		entry.Keys = keys
		entry.Status = "success"
		_ = s.logger.Log(entry)
	}
	return nil
}

// Ensure vault.Client satisfies VaultReader.
var _ VaultReader = (*vault.Client)(nil)
