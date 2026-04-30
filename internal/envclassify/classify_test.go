package envclassify_test

import (
	"testing"

	"github.com/nahK994/vaultpull/internal/envclassify"
)

func baseRules() []envclassify.Rule {
	return []envclassify.Rule{
		{Pattern: `(?i)(secret|token|password|key)`, Category: envclassify.CategorySecret},
		{Pattern: `(?i)^DB_`, Category: envclassify.CategoryDatabase},
		{Pattern: `(?i)^FEATURE_`, Category: envclassify.CategoryFeature},
		{Pattern: `(?i)(host|port|url|addr)`, Category: envclassify.CategoryConfig},
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := envclassify.New(baseRules())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	rules := []envclassify.Rule{{Pattern: "[invalid", Category: envclassify.CategorySecret}}
	_, err := envclassify.New(rules)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	rules := []envclassify.Rule{{Pattern: "", Category: envclassify.CategoryConfig}}
	_, err := envclassify.New(rules)
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestClassify_SecretKey(t *testing.T) {
	c, _ := envclassify.New(baseRules())
	secrets := map[string]string{"API_TOKEN": "abc123"}
	results := c.Classify(secrets)
	if len(results) != 1 || results[0].Category != envclassify.CategorySecret {
		t.Fatalf("expected secret, got %v", results)
	}
}

func TestClassify_DatabaseKey(t *testing.T) {
	c, _ := envclassify.New(baseRules())
	secrets := map[string]string{"DB_HOST": "localhost"}
	results := c.Classify(secrets)
	if len(results) != 1 || results[0].Category != envclassify.CategoryDatabase {
		t.Fatalf("expected database, got %v", results)
	}
}

func TestClassify_UnknownKey(t *testing.T) {
	c, _ := envclassify.New(baseRules())
	secrets := map[string]string{"FOOBAR": "baz"}
	results := c.Classify(secrets)
	if len(results) != 1 || results[0].Category != envclassify.CategoryUnknown {
		t.Fatalf("expected unknown, got %v", results)
	}
}

func TestSummary_CountsCategories(t *testing.T) {
	c, _ := envclassify.New(baseRules())
	secrets := map[string]string{
		"API_SECRET": "x",
		"DB_NAME":    "mydb",
		"FOOBAR":     "baz",
	}
	results := c.Classify(secrets)
	summary := envclassify.Summary(results)
	if summary[envclassify.CategorySecret] != 1 {
		t.Errorf("expected 1 secret, got %d", summary[envclassify.CategorySecret])
	}
	if summary[envclassify.CategoryDatabase] != 1 {
		t.Errorf("expected 1 database, got %d", summary[envclassify.CategoryDatabase])
	}
	if summary[envclassify.CategoryUnknown] != 1 {
		t.Errorf("expected 1 unknown, got %d", summary[envclassify.CategoryUnknown])
	}
}
