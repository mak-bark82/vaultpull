package envmask

import (
	"strings"
)

// MaskMode controls how values are masked.
type MaskMode int

const (
	// MaskFull replaces the entire value with asterisks.
	MaskFull MaskMode = iota
	// MaskPartial reveals the last N characters.
	MaskPartial
)

// Options configures masking behaviour.
type Options struct {
	Mode        MaskMode
	VisibleChars int   // used when Mode == MaskPartial
	MaskChar    rune
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Mode:        MaskPartial,
		VisibleChars: 4,
		MaskChar:    '*',
	}
}

// Masker applies masking rules to a set of env key/value pairs.
type Masker struct {
	opts    Options
	patterns []string
}

// New creates a Masker with the given options and key patterns.
// Patterns are matched as case-insensitive substrings against key names.
func New(opts Options, patterns []string) *Masker {
	lower := make([]string, len(patterns))
	for i, p := range patterns {
		lower[i] = strings.ToLower(p)
	}
	return &Masker{opts: opts, patterns: lower}
}

// Apply returns a copy of secrets with sensitive values masked.
func (m *Masker) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if m.isSensitive(k) {
			out[k] = m.mask(v)
		} else {
			out[k] = v
		}
	}
	return out
}

func (m *Masker) isSensitive(key string) bool {
	lk := strings.ToLower(key)
	for _, p := range m.patterns {
		if strings.Contains(lk, p) {
			return true
		}
	}
	return false
}

func (m *Masker) mask(value string) string {
	if m.opts.Mode == MaskFull || len(value) == 0 {
		return strings.Repeat(string(m.opts.MaskChar), len(value))
	}
	// MaskPartial: show last VisibleChars characters
	visible := m.opts.VisibleChars
	if visible >= len(value) {
		return value
	}
	hidden := len(value) - visible
	return strings.Repeat(string(m.opts.MaskChar), hidden) + value[hidden:]
}
