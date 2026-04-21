// Package envhash provides utilities for computing and comparing
// stable hashes of environment variable maps, useful for detecting
// changes between syncs without exposing secret values.
package envhash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Hasher computes deterministic hashes over env maps.
type Hasher struct{}

// New returns a new Hasher.
func New() *Hasher {
	return &Hasher{}
}

// Hash returns a SHA-256 hex digest of the sorted key=value pairs
// in the provided map. The result is stable regardless of map
// iteration order.
func (h *Hasher) Hash(env map[string]string) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
	}

	sum := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(sum[:])
}

// HashKey returns a SHA-256 hex digest of a single key's value,
// useful for detecting per-key rotation without revealing the value.
func (h *Hasher) HashKey(key, value string) string {
	sum := sha256.Sum256([]byte(key + "=" + value))
	return hex.EncodeToString(sum[:])
}

// Equal reports whether two env maps produce the same hash.
func (h *Hasher) Equal(a, b map[string]string) bool {
	return h.Hash(a) == h.Hash(b)
}

// Diff returns the set of keys whose hashed values differ between
// old and new, including keys that were added or removed.
func (h *Hasher) Diff(old, next map[string]string) []string {
	seen := map[string]struct{}{}
	var changed []string

	for k, v := range next {
		seen[k] = struct{}{}
		if ov, ok := old[k]; !ok || h.HashKey(k, v) != h.HashKey(k, ov) {
			changed = append(changed, k)
		}
	}
	for k := range old {
		if _, ok := seen[k]; !ok {
			changed = append(changed, k)
		}
	}
	sort.Strings(changed)
	return changed
}
