package envclone

import (
	"errors"
	"fmt"
)

// CloneOptions controls how secrets are cloned between vault paths.
type CloneOptions struct {
	Overwrite bool
	DryRun    bool
}

// Result holds the outcome of a clone operation.
type Result struct {
	Copied    []string
	Skipped   []string
	Overwrite []string
}

// Summary returns a human-readable summary of the clone result.
func (r *Result) Summary() string {
	return fmt.Sprintf("copied=%d skipped=%d overwritten=%d",
		len(r.Copied), len(r.Skipped), len(r.Overwrite))
}

// Clone copies secrets from src into dst according to the given options.
// Keys already present in dst are skipped unless Overwrite is true.
// When DryRun is true, dst is never modified.
func Clone(src, dst map[string]string, opts CloneOptions) (*Result, error) {
	if src == nil {
		return nil, errors.New("envclone: source map must not be nil")
	}
	if dst == nil {
		return nil, errors.New("envclone: destination map must not be nil")
	}

	result := &Result{}

	for k, v := range src {
		_, exists := dst[k]
		switch {
		case !exists:
			if !opts.DryRun {
				dst[k] = v
			}
			result.Copied = append(result.Copied, k)
		case exists && opts.Overwrite:
			if !opts.DryRun {
				dst[k] = v
			}
			result.Overwrite = append(result.Overwrite, k)
		default:
			result.Skipped = append(result.Skipped, k)
		}
	}

	return result, nil
}
