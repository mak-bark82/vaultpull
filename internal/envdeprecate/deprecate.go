package envdeprecate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a deprecation rule for an environment variable key.
type Rule struct {
	Key        string `yaml:"key"`
	Pattern    string `yaml:"pattern"`
	Message    string `yaml:"message"`
	Replacement string `yaml:"replacement"`
	compiled   *regexp.Regexp
}

// Finding represents a single deprecation match.
type Finding struct {
	Key         string
	Message     string
	Replacement string
}

// Checker evaluates env maps against deprecation rules.
type Checker struct {
	rules []Rule
}

// New creates a Checker from a set of rules.
// Returns an error if any rule has an invalid pattern.
func New(rules []Rule) (*Checker, error) {
	compiled := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if r.Key == "" && r.Pattern == "" {
			return nil, fmt.Errorf("rule must specify key or pattern")
		}
		if r.Pattern != "" {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				return nil, fmt.Errorf("invalid pattern %q: %w", r.Pattern, err)
			}
			r.compiled = re
		}
		if r.Message == "" {
			r.Message = fmt.Sprintf("%s is deprecated", r.Key)
		}
		compiled = append(compiled, r)
	}
	return &Checker{rules: compiled}, nil
}

// Check returns a list of findings for any deprecated keys found in env.
func (c *Checker) Check(env map[string]string) []Finding {
	var findings []Finding
	for key := range env {
		for _, r := range c.rules {
			if matches(r, key) {
				findings = append(findings, Finding{
					Key:         key,
					Message:     r.Message,
					Replacement: r.Replacement,
				})
				break
			}
		}
	}
	return findings
}

// Summary returns a human-readable summary of findings.
func Summary(findings []Finding) string {
	if len(findings) == 0 {
		return "no deprecated keys found"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d deprecated key(s) found:\n", len(findings))
	for _, f := range findings {
		line := fmt.Sprintf("  - %s: %s", f.Key, f.Message)
		if f.Replacement != "" {
			line += fmt.Sprintf(" (use %s instead)", f.Replacement)
		}
		sb.WriteString(line + "\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}

func matches(r Rule, key string) bool {
	if r.Key != "" && r.Key == key {
		return true
	}
	if r.compiled != nil && r.compiled.MatchString(key) {
		return true
	}
	return false
}
