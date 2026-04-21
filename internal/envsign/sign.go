// Package envSign provides HMAC-based signing and verification of env secret maps.
// It allows detecting tampering by signing the canonical representation of secrets.
package envSign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Signer signs and verifies env secret maps using HMAC-SHA256.
type Signer struct {
	key []byte
}

// New creates a new Signer with the given secret key.
// Returns an error if the key is empty.
func New(key []byte) (*Signer, error) {
	if len(key) == 0 {
		return nil, errors.New("envSign: key must not be empty")
	}
	return &Signer{key: key}, nil
}

// Sign computes an HMAC-SHA256 signature over the canonical form of secrets.
// The canonical form is: sorted KEY=VALUE lines joined by newlines.
func (s *Signer) Sign(secrets map[string]string) string {
	mac := hmac.New(sha256.New, s.key)
	mac.Write([]byte(canonical(secrets)))
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify checks whether the given signature matches the secrets map.
// Returns nil if valid, or a descriptive error if tampered or mismatched.
func (s *Signer) Verify(secrets map[string]string, sig string) error {
	expected := s.Sign(secrets)
	if !hmac.Equal([]byte(expected), []byte(sig)) {
		return fmt.Errorf("envSign: signature mismatch: verification failed")
	}
	return nil
}

// canonical returns a deterministic string representation of the secrets map.
func canonical(secrets map[string]string) string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys))
	for _, k := range keys {
		lines = append(lines, k+"="+secrets[k])
	}
	return strings.Join(lines, "\n")
}
