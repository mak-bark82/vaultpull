package envobfuscate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a pattern and the obfuscation strategy to apply.
type Rule struct {
	Pattern     string `yaml:"pattern"`
	Strategy    string `yaml:"strategy"` // "hash", "mask", "remove"
	RevealChars int    `yaml:"reveal_chars"`
}

// Result holds the obfuscated environment map and a summary of changes.
type Result struct {
	Env     map[string]string
	Changed []string
}

// Obfuscator applies obfuscation rules to environment variable values.
type Obfuscator struct {
	rules   []compiledRule
}

type compiledRule struct {
	Rule
	re *regexp.Regexp
}

// New creates an Obfuscator from the given rules.
// Returns an error if any pattern fails to compile.
func New(rules []Rule) (*Obfuscator, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("envobfuscate: at least one rule is required")
	}
	compiled := make([]compiledRule, 0, len(rules))
	for _, r := range rules {
		if strings.TrimSpace(r.Pattern) == "" {
			return nil, fmt.Errorf("envobfuscate: rule pattern must not be empty")
		}
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, fmt.Errorf("envobfuscate: invalid pattern %q: %w", r.Pattern, err)
		}
		compiled = append(compiled, compiledRule{Rule: r, re: re})
	}
	return &Obfuscator{rules: compiled}, nil
}

// Apply obfuscates matching keys in the provided env map.
// The original map is not mutated.
func (o *Obfuscator) Apply(env map[string]string) Result {
	out := make(map[string]string, len(env))
	var changed []string
	for k, v := range env {
		out[k] = v
	}
	for k, v := range env {
		for _, r := range o.rules {
			if r.re.MatchString(k) {
				obfuscated := applyStrategy(v, r.Rule)
				if obfuscated != v {
					out[k] = obfuscated
					changed = append(changed, k)
				}
				break
			}
		}
	}
	return Result{Env: out, Changed: changed}
}

func applyStrategy(value string, r Rule) string {
	switch r.Strategy {
	case "remove":
		return ""
	case "hash":
		return fmt.Sprintf("[redacted:%d]", len(value))
	case "mask":
		if r.RevealChars <= 0 || len(value) <= r.RevealChars {
			return strings.Repeat("*", len(value))
		}
		return strings.Repeat("*", len(value)-r.RevealChars) + value[len(value)-r.RevealChars:]
	default:
		return strings.Repeat("*", len(value))
	}
}
