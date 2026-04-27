package envencrypt

import (
	"errors"
	"fmt"
	"regexp"
)

// Rule defines an encryption rule: keys matching Pattern will have their
// values encrypted or decrypted using the associated Cipher.
type Rule struct {
	Pattern string
	regexp  *regexp.Regexp
}

// Processor applies encryption or decryption rules to env maps.
type Processor struct {
	rules  []Rule
	cipher Cipher
}

// Cipher is the interface required by Processor for encrypt/decrypt operations.
type Cipher interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// New creates a Processor with the given cipher and key patterns.
// Patterns are compiled as regular expressions.
func New(c Cipher, patterns []string) (*Processor, error) {
	if c == nil {
		return nil, errors.New("envencrypt: cipher must not be nil")
	}
	rules := make([]Rule, 0, len(patterns))
	for _, p := range patterns {
		if p == "" {
			return nil, errors.New("envencrypt: pattern must not be empty")
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("envencrypt: invalid pattern %q: %w", p, err)
		}
		rules = append(rules, Rule{Pattern: p, regexp: re})
	}
	return &Processor{rules: rules, cipher: c}, nil
}

// Encrypt returns a new map where values whose keys match any rule pattern
// are replaced with their encrypted form.
func (p *Processor) Encrypt(env map[string]string) (map[string]string, error) {
	return p.process(env, true)
}

// Decrypt returns a new map where values whose keys match any rule pattern
// are replaced with their decrypted form.
func (p *Processor) Decrypt(env map[string]string) (map[string]string, error) {
	return p.process(env, false)
}

func (p *Processor) process(env map[string]string, encrypt bool) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if p.matches(k) {
			var err error
			var transformed string
			if encrypt {
				transformed, err = p.cipher.Encrypt(v)
			} else {
				transformed, err = p.cipher.Decrypt(v)
			}
			if err != nil {
				return nil, fmt.Errorf("envencrypt: key %q: %w", k, err)
			}
			out[k] = transformed
		} else {
			out[k] = v
		}
	}
	return out, nil
}

func (p *Processor) matches(key string) bool {
	for _, r := range p.rules {
		if r.regexp.MatchString(key) {
			return true
		}
	}
	return false
}
