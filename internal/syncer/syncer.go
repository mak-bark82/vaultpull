package syncer

import (
	"fmt"

	"github.com/user/vaultpull/internal/audit"
	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/diff"
	"github.com/user/vaultpull/internal/envreader"
	"github.com/user/vaultpull/internal/envwriter"
	"github.com/user/vaultpull/internal/vault"
)

// Syncer orchestrates reading secrets from Vault and writing them to .env files.
type Syncer struct {
	client   *vault.Client
	mappings []config.Mapping
	logger   *audit.Logger
}

// New creates a Syncer from the given config.
func New(cfg *config.Config, mappings []config.Mapping, logger *audit.Logger) (*Syncer, error) {
	c, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return nil, err
	}
	return &Syncer{client: c, mappings: mappings, logger: logger}, nil
}

// Run iterates over all mappings, diffs secrets, and writes changes.
func (s *Syncer) Run() error {
	for _, m := range s.mappings {
		secrets, err := s.client.ReadSecrets(m.VaultPath)
		if err != nil {
			return fmt.Errorf("reading vault path %q: %w", m.VaultPath, err)
		}

		existing, err := envreader.Read(m.EnvFile)
		if err != nil {
			return fmt.Errorf("reading env file %q: %w", m.EnvFile, err)
		}

		result := diff.Compare(existing, secrets)
		if !result.HasChanges() {
			fmt.Printf("%s: no changes\n", m.EnvFile)
			continue
		}

		if err := envwriter.Write(m.EnvFile, secrets, true); err != nil {
			return fmt.Errorf("writing env file %q: %w", m.EnvFile, err)
		}

		_ = s.logger.Log(m.VaultPath, m.EnvFile, result.Summary())
		fmt.Printf("%s: %s\n", m.EnvFile, result.Summary())
	}
	return nil
}
