// Package envresolve resolves the final set of secrets to write by applying
// profile, scope, filter, transform, and expand stages in sequence.
package envresolve

import (
	"fmt"

	"github.com/your-org/vaultpull/internal/envexpand"
	"github.com/your-org/vaultpull/internal/envfilter"
	"github.com/your-org/vaultpull/internal/envtransform"
)

// Options controls which resolution stages are active.
type Options struct {
	// IncludeKeys limits output to these keys (empty = all).
	IncludeKeys []string
	// ExcludeKeys removes these keys from output.
	ExcludeKeys []string
	// Transforms is a map of key → TransformRule applied after filtering.
	Transforms map[string]envtransform.Rule
	// ExpandRefs resolves ${VAR} references when true.
	ExpandRefs bool
}

// Resolver applies a deterministic pipeline to a raw secrets map.
type Resolver struct {
	opts Options
}

// New returns a Resolver configured with opts.
func New(opts Options) *Resolver {
	return &Resolver{opts: opts}
}

// Resolve runs the full pipeline: filter → transform → expand.
// It returns a new map and never mutates secrets.
func (r *Resolver) Resolve(secrets map[string]string) (map[string]string, error) {
	if secrets == nil {
		return map[string]string{}, nil
	}

	// Stage 1 – filter
	filtered := envfilter.Filter(secrets, envfilter.Rules{
		Include: r.opts.IncludeKeys,
		Exclude: r.opts.ExcludeKeys,
	})

	// Stage 2 – transform
	var transformed map[string]string
	if len(r.opts.Transforms) > 0 {
		var err error
		transformed, err = envtransform.Transform(filtered, r.opts.Transforms)
		if err != nil {
			return nil, fmt.Errorf("envresolve: transform stage: %w", err)
		}
	} else {
		transformed = filtered
	}

	// Stage 3 – expand
	if r.opts.ExpandRefs {
		return envexpand.Expand(transformed), nil
	}
	return transformed, nil
}
