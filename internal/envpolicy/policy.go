package envpolicy

import (
	"fmt"
	"regexp"
	"strings"
)

// Action defines what to do when a policy rule matches.
type Action string

const (
	ActionAllow Action = "allow"
	ActionDeny  Action = "deny"
	ActionWarn  Action = "warn"
)

// Rule represents a single policy rule applied to env keys or values.
type Rule struct {
	Name    string `yaml:"name"`
	Pattern string `yaml:"pattern"`
	Target  string `yaml:"target"` // "key" or "value"
	Action  Action `yaml:"action"`

	re *regexp.Regexp
}

// Violation is returned when a rule is triggered.
type Violation struct {
	Rule    string
	Key     string
	Action  Action
	Message string
}

// Enforcer evaluates a set of rules against env secrets.
type Enforcer struct {
	rules []Rule
}

// New creates an Enforcer from a slice of rules, compiling regex patterns.
func New(rules []Rule) (*Enforcer, error) {
	for i, r := range rules {
		if r.Pattern == "" {
			return nil, fmt.Errorf("rule %q has empty pattern", r.Name)
		}
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, fmt.Errorf("rule %q: invalid pattern: %w", r.Name, err)
		}
		rules[i].re = re
	}
	return &Enforcer{rules: rules}, nil
}

// Check evaluates all rules against the provided secrets map.
// Returns a list of violations (may include warns and denies).
func (e *Enforcer) Check(secrets map[string]string) []Violation {
	var violations []Violation
	for key, val := range secrets {
		for _, rule := range e.rules {
			var subject string
			switch strings.ToLower(rule.Target) {
			case "value":
				subject = val
			default:
				subject = key
			}
			if rule.re.MatchString(subject) {
				violations = append(violations, Violation{
					Rule:    rule.Name,
					Key:     key,
					Action:  rule.Action,
					Message: fmt.Sprintf("rule %q matched %s %q", rule.Name, rule.Target, subject),
				})
			}
		}
	}
	return violations
}

// HasDenials returns true if any violation carries ActionDeny.
func HasDenials(violations []Violation) bool {
	for _, v := range violations {
		if v.Action == ActionDeny {
			return true
		}
	}
	return false
}
