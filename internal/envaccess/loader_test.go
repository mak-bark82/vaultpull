package envaccess_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/envaccess"
)

func writeRulesFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "access.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	return p
}

func TestLoadRules_Valid(t *testing.T) {
	p := writeRulesFile(t, "rules:\n  - pattern: DB_*\n    permission: write\n  - pattern: API_KEY\n    permission: read\n")
	rules, err := envaccess.LoadRules(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Permission != envaccess.PermWrite {
		t.Errorf("expected PermWrite for DB_*")
	}
	if rules[1].Permission != envaccess.PermRead {
		t.Errorf("expected PermRead for API_KEY")
	}
}

func TestLoadRules_EmptyPath(t *testing.T) {
	rules, err := envaccess.LoadRules("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected empty rules, got %d", len(rules))
	}
}

func TestLoadRules_MissingFile(t *testing.T) {
	_, err := envaccess.LoadRules("/nonexistent/access.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadRules_UnknownPermission(t *testing.T) {
	p := writeRulesFile(t, "rules:\n  - pattern: FOO\n    permission: superadmin\n")
	_, err := envaccess.LoadRules(p)
	if err == nil {
		t.Fatal("expected error for unknown permission")
	}
}

func TestLoadRules_NonePermission(t *testing.T) {
	p := writeRulesFile(t, "rules:\n  - pattern: SECRET_*\n    permission: none\n")
	rules, err := envaccess.LoadRules(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules[0].Permission != envaccess.PermNone {
		t.Errorf("expected PermNone")
	}
}
