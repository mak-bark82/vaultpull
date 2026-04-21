// Package envsync orchestrates a full end-to-end sync cycle: fetch secrets
// from Vault, apply filters and transforms, diff against the existing .env
// file, optionally back it up, then write the merged result.
package envsync

import (
	"fmt"
	"io"

	"github.com/yourusername/vaultpull/internal/envbackup"
	"github.com/yourusername/vaultpull/internal/envfilter"
	"github.com/yourusername/vaultpull/internal/envmerge"
	"github.com/yourusername/vaultpull/internal/envreader"
	"github.com/yourusername/vaultpull/internal/envtransform"
	"github.com/yourusername/vaultpull/internal/envwriter"
	"github.com/yourusername/vaultpull/internal/diff"
)

// SecretFetcher retrieves key/value secrets from a remote source (e.g. Vault).
type SecretFetcher interface {
	ReadSecrets(path string) (map[string]string, error)
}

// Options controls the behaviour of a sync run.
type Options struct {
	// VaultPath is the KV path to read from Vault.
	VaultPath string

	// EnvFile is the local .env file to write results into.
	EnvFile string

	// Overwrite controls whether existing keys are replaced by Vault values.
	Overwrite bool

	// Backup creates a timestamped copy of the existing file before writing.
	Backup bool

	// Filter is applied to the raw secrets fetched from Vault.
	Filter envfilter.Rules

	// Transform is applied after filtering.
	Transform envtransform.Rules

	// Out receives a human-readable summary; may be nil to suppress output.
	Out io.Writer
}

// Result summarises what happened during a sync.
type Result struct {
	Added   int
	Changed int
	Removed int
	Total   int
}

// String returns a one-line summary suitable for CLI output.
func (r Result) String() string {
	return fmt.Sprintf("synced: +%d ~%d -%d (total %d)", r.Added, r.Changed, r.Removed, r.Total)
}

// Run executes a full sync cycle using the provided fetcher and options.
// It is the primary entry-point for callers that want programmatic control.
func Run(fetcher SecretFetcher, opts Options) (Result, error) {
	// 1. Fetch secrets from Vault.
	vaultSecrets, err := fetcher.ReadSecrets(opts.VaultPath)
	if err != nil {
		return Result{}, fmt.Errorf("envsync: fetch %q: %w", opts.VaultPath, err)
	}

	// 2. Apply filter rules.
	filtered := envfilter.Filter(vaultSecrets, opts.Filter)

	// 3. Apply transform rules.
	transformed := envtransform.Transform(filtered, opts.Transform)

	// 4. Read the existing local .env (empty map if file does not exist yet).
	existing, err := envreader.Read(opts.EnvFile)
	if err != nil {
		// Non-existent file is acceptable; treat as empty.
		existing = map[string]string{}
	}

	// 5. Diff so we can report what will change.
	changes := diff.Compare(existing, transformed)

	// 6. Optionally back up the current file before mutating it.
	if opts.Backup {
		if _, berr := envbackup.Backup(opts.EnvFile); berr != nil {
			return Result{}, fmt.Errorf("envsync: backup %q: %w", opts.EnvFile, berr)
		}
	}

	// 7. Merge Vault values into the existing map.
	merged, mergeStats := envmerge.Merge(existing, transformed, opts.Overwrite)

	// 8. Write the merged result back to disk.
	if werr := envwriter.Write(opts.EnvFile, merged); werr != nil {
		return Result{}, fmt.Errorf("envsync: write %q: %w", opts.EnvFile, werr)
	}

	// 9. Build result summary.
	res := Result{
		Added:   changes.Added,
		Changed: changes.Changed,
		Removed: changes.Removed,
		Total:   mergeStats.Total,
	}

	// 10. Emit summary to writer when provided.
	if opts.Out != nil {
		fmt.Fprintln(opts.Out, res.String())
	}

	return res, nil
}
