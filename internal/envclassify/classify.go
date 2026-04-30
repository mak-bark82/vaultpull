package envclassify

import (
	"regexp"
	"strings"
)

// Category represents a classification label for an env key.
type Category string

const (
	CategorySecret   Category = "secret"
	CategoryConfig   Category = "config"
	CategoryFeature  Category = "feature"
	CategoryDatabase Category = "database"
	CategoryUnknown  Category = "unknown"
)

// Rule maps a regexp pattern to a category.
type Rule struct {
	Pattern  string   `yaml:"pattern"`
	Category Category `yaml:"category"`
	re       *regexp.Regexp
}

// Result holds the classification result for a single key.
type Result struct {
	Key      string
	Category Category
}

// Classifier assigns categories to env keys based on rules.
type Classifier struct {
	rules []Rule
}

// New creates a Classifier from a slice of rules.
// Rules with invalid patterns are skipped.
func New(rules []Rule) (*Classifier, error) {
	compiled := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if strings.TrimSpace(r.Pattern) == "" {
			return nil, fmt.Errorf("envclassify: empty pattern in rule")
		}
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, fmt.Errorf("envclassify: invalid pattern %q: %w", r.Pattern, err)
		}
		r.re = re
		compiled = append(compiled, r)
	}
	return &Classifier{rules: compiled}, nil
}

// Classify returns a Result for each key in secrets.
// Keys that match no rule are assigned CategoryUnknown.
func (c *Classifier) Classify(secrets map[string]string) []Result {
	results := make([]Result, 0, len(secrets))
	for k := range secrets {
		cat := CategoryUnknown
		for _, r := range c.rules {
			if r.re.MatchString(k) {
				cat = r.Category
				break
			}
		}
		results = append(results, Result{Key: k, Category: cat})
	}
	return results
}

// Summary returns a count of keys per category.
func Summary(results []Result) map[Category]int {
	m := make(map[Category]int)
	for _, r := range results {
		m[r.Category]++
	}
	return m
}
