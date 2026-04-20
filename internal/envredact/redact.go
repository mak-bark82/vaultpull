package envredact

import "strings"

// DefaultPatterns are key substrings that trigger redaction.
var DefaultPatterns = []string{
	"password", "secret", "token", "key", "api", "auth", "credential",
}

// Redactor masks sensitive secret values before display or logging.
type Redactor struct {
	patterns []string
	mask     string
}

// New creates a Redactor with the given patterns and mask string.
// If patterns is nil, DefaultPatterns are used. If mask is empty, "***" is used.
func New(patterns []string, mask string) *Redactor {
	if patterns == nil {
		patterns = DefaultPatterns
	}
	if mask == "" {
		mask = "***"
	}
	return &Redactor{patterns: patterns, mask: mask}
}

// IsSensitive reports whether the given key matches any redaction pattern.
func (r *Redactor) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range r.patterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// Redact returns a copy of secrets where sensitive values are replaced with the mask.
func (r *Redactor) Redact(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if r.IsSensitive(k) {
			out[k] = r.mask
		} else {
			out[k] = v
		}
	}
	return out
}
