package envaccess

import (
	"errors"
	"fmt"	
	"strings"
)

// Permission represents an access level for a secret key.
type Permission int

const (
	PermNone  Permission = iota
	PermRead             // read-only
	PermWrite            // read + write
)

// Rule defines access control for a key pattern.
type Rule struct {
	Pattern    string     `yaml:"pattern"`
	Permission Permission `yaml:"permission"`
}

// Checker enforces access rules against secret keys.
type Checker struct {
	rules []Rule
}

// New creates a Checker from the given rules.
func New(rules []Rule) (*Checker, error) {
	for _, r := range rules {
		if strings.TrimSpace(r.Pattern) == "" {
			return nil, errors.New("envaccess: rule has empty pattern")
		}
	}
	return &Checker{rules: rules}, nil
}

// Check returns the effective Permission for the given key.
// Rules are evaluated in order; the first match wins.
// If no rule matches, PermNone is returned.
func (c *Checker) Check(key string) Permission {
	for _, r := range c.rules {
		if matchPattern(r.Pattern, key) {
			return r.Permission
		}
	}
	return PermNone
}

// Enforce returns an error if the key does not satisfy the required permission.
func (c *Checker) Enforce(key string, required Permission) error {
	got := c.Check(key)
	if got < required {
		return fmt.Errorf("envaccess: key %q has permission %d, need %d", key, got, required)
	}
	return nil
}

// matchPattern supports a trailing '*' wildcard.
func matchPattern(pattern, key string) bool {
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(key, pattern[:len(pattern)-1])
	}
	return pattern == key
}
