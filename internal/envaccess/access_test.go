package envaccess_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/envaccess"
)

func baseRules() []envaccess.Rule {
	return []envaccess.Rule{
		{Pattern: "DB_*", Permission: envaccess.PermWrite},
		{Pattern: "API_KEY", Permission: envaccess.PermRead},
		{Pattern: "SECRET_*", Permission: envaccess.PermNone},
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := envaccess.New(baseRules())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	rules := []envaccess.Rule{{Pattern: "", Permission: envaccess.PermRead}}
	_, err := envaccess.New(rules)
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestCheck_WildcardMatch(t *testing.T) {
	c, _ := envaccess.New(baseRules())
	if got := c.Check("DB_HOST"); got != envaccess.PermWrite {
		t.Errorf("expected PermWrite, got %d", got)
	}
}

func TestCheck_ExactMatch(t *testing.T) {
	c, _ := envaccess.New(baseRules())
	if got := c.Check("API_KEY"); got != envaccess.PermRead {
		t.Errorf("expected PermRead, got %d", got)
	}
}

func TestCheck_NoMatch_ReturnsNone(t *testing.T) {
	c, _ := envaccess.New(baseRules())
	if got := c.Check("UNKNOWN_VAR"); got != envaccess.PermNone {
		t.Errorf("expected PermNone, got %d", got)
	}
}

func TestCheck_FirstRuleWins(t *testing.T) {
	rules := []envaccess.Rule{
		{Pattern: "DB_*", Permission: envaccess.PermRead},
		{Pattern: "DB_HOST", Permission: envaccess.PermWrite},
	}
	c, _ := envaccess.New(rules)
	if got := c.Check("DB_HOST"); got != envaccess.PermRead {
		t.Errorf("expected PermRead (first rule wins), got %d", got)
	}
}

func TestEnforce_Allowed(t *testing.T) {
	c, _ := envaccess.New(baseRules())
	if err := c.Enforce("DB_PASS", envaccess.PermRead); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestEnforce_Denied(t *testing.T) {
	c, _ := envaccess.New(baseRules())
	if err := c.Enforce("API_KEY", envaccess.PermWrite); err == nil {
		t.Fatal("expected error for insufficient permission")
	}
}
