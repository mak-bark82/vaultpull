package syncer

import (
	"fmt"

	"github.com/yourusername/vaultpull/internal/audit"
	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/dryrun"
	"github.com/yourusername/vaultpull/internal/envreader"
	"github.com/yourusername/vaultpull/internal/envwriter"
	"github.com/yourusername/vaultpull/internal/vault"
)

// Syncer orchestrates pulling secrets from Vault and writing them to env files.
type Syncer struct {
	client   *vault.Client
	mappings []config.Mapping
	logger   *audit.Logger
	dryRun   bool
	reporter *dryrun.Reporter
}

// New creates a Syncer from the provided config.
func New(cfg *config.Config, mappings []config.Mapping, logger *audit.Logger) (*Syncer, error) {
	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return nil, fmt.Errorf("vault client: %w", err)
	}
	return &Syncer{
		client:   client,
		mappings: mappings,
		logger:   logger,
		dryRun:   cfg.DryRun,
		reporter: dryrun.NewReporter(nil),
	}, nil
}

// Run iterates over mappings, reads secrets from Vault, and writes env files.
func (s *Syncer) Run() error {
	for _, m := range s.mappings {
		secrets, err := s.client.ReadSecrets(m.VaultPath)
		if err != nil {
			return fmt.Errorf("read %s: %w", m.VaultPath, err)
		}
		existing, _ := envreader.Read(m.EnvFile)
		if s.dryRun {
			s.reporter.Report(m.EnvFile, existing, secrets)
			continue
		}
		if err := envwriter.Write(m.EnvFile, secrets, m.Overwrite); err != nil {
			return fmt.Errorf("write %s: %w", m.EnvFile, err)
		}
		if s.logger != nil {
			_ = s.logger.Log(m.VaultPath, m.EnvFile, secrets)
		}
	}
	return nil
}
