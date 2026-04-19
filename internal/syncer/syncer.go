package syncer

import (
	"fmt"

	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/envwriter"
	"github.com/example/vaultpull/internal/vault"
)

// Result holds the outcome of a single secret sync operation.
type Result struct {
	Path    string
	OutFile string
	Skipped bool
	Err     error
}

// Syncer orchestrates reading secrets from Vault and writing them to env files.
type Syncer struct {
	client *vault.Client
	cfg    *config.Config
}

// New creates a new Syncer from the provided config.
func New(cfg *config.Config) (*Syncer, error) {
	client, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return nil, fmt.Errorf("syncer: failed to create vault client: %w", err)
	}
	return &Syncer{client: client, cfg: cfg}, nil
}

// Run iterates over all configured secret mappings and syncs each one.
func (s *Syncer) Run() []Result {
	results := make([]Result, 0, len(s.cfg.Mappings))
	for _, m := range s.cfg.Mappings {
		res := s.syncOne(m)
		results = append(results, res)
	}
	return results
}

func (s *Syncer) syncOne(m config.Mapping) Result {
	res := Result{Path: m.VaultPath, OutFile: m.EnvFile}

	secrets, err := s.client.ReadSecrets(m.VaultPath)
	if err != nil {
		res.Err = fmt.Errorf("reading %s: %w", m.VaultPath, err)
		return res
	}

	skipped, err := envwriter.Write(m.EnvFile, secrets, m.Overwrite)
	if err != nil {
		res.Err = fmt.Errorf("writing %s: %w", m.EnvFile, err)
		return res
	}
	res.Skipped = skipped
	return res
}
