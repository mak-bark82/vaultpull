package envdeprecate_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/envdeprecate"
)

func baseEnv() map[string]string {
	return map[string]string{
		"OLD_API_KEY":  "abc",
		"LEGACY_TOKEN": "xyz",
		"NEW_SECRET":   "ok",
	}
}

func TestCheck_NoFindings(t *testing.T) {
	c, err := envdeprecate.New([]envdeprecate.Rule{
		{Key: "REMOVED_KEY", Message: "gone"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	findings := c.Check(baseEnv())
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d", len(findings))
	}
}

func TestCheck_ExactKeyMatch(t *testing.T) {
	c, err := envdeprecate.New([]envdeprecate.Rule{
		{Key: "OLD_API_KEY", Message: "use NEW_API_KEY", Replacement: "NEW_API_KEY"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	findings := c.Check(baseEnv())
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "OLD_API_KEY" {
		t.Errorf("unexpected key: %s", findings[0].Key)
	}
	if findings[0].Replacement != "NEW_API_KEY" {
		t.Errorf("unexpected replacement: %s", findings[0].Replacement)
	}
}

func TestCheck_PatternMatch(t *testing.T) {
	c, err := envdeprecate.New([]envdeprecate.Rule{
		{Pattern: "^LEGACY_", Message: "legacy keys are deprecated"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	findings := c.Check(baseEnv())
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "LEGACY_TOKEN" {
		t.Errorf("unexpected key: %s", findings[0].Key)
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := envdeprecate.New([]envdeprecate.Rule{
		{Pattern: "[invalid"},
	})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_MissingKeyAndPattern(t *testing.T) {
	_, err := envdeprecate.New([]envdeprecate.Rule{
		{Message: "no key or pattern"},
	})
	if err == nil {
		t.Fatal("expected error when key and pattern are both empty")
	}
}

func TestSummary_NoFindings(t *testing.T) {
	s := envdeprecate.Summary(nil)
	if s != "no deprecated keys found" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_WithFindings(t *testing.T) {
	findings := []envdeprecate.Finding{
		{Key: "OLD_KEY", Message: "deprecated", Replacement: "NEW_KEY"},
	}
	s := envdeprecate.Summary(findings)
	if !strings.Contains(s, "OLD_KEY") {
		t.Errorf("summary missing key: %s", s)
	}
	if !strings.Contains(s, "NEW_KEY") {
		t.Errorf("summary missing replacement: %s", s)
	}
}

func TestLoadRules_Valid(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "rules.yaml")
	content := `rules:
  - key: OLD_TOKEN
    message: use NEW_TOKEN
    replacement: NEW_TOKEN
  - pattern: "^DEPRECATED_"
    message: all DEPRECATED_ keys are removed
`
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	rules, err := envdeprecate.LoadRules(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
}

func TestLoadRules_EmptyPath(t *testing.T) {
	rules, err := envdeprecate.LoadRules("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected empty rules")
	}
}

func TestLoadRules_MissingFile(t *testing.T) {
	_, err := envdeprecate.LoadRules("/nonexistent/rules.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
