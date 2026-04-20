package envscope_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envscope"
)

func writeScopeFile(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "scopes.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write scope file: %v", err)
	}
	return p
}

func TestLoadScopes_Valid(t *testing.T) {
	p := writeScopeFile(t, "scopes:\n  - name: dev\n    prefix: secret/dev\n  - name: prod\n    prefix: secret/prod\n")
	r, err := envscope.LoadScopes(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := r.Resolve("prod", "app")
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}
	if got != "secret/prod/app" {
		t.Errorf("unexpected path: %s", got)
	}
}

func TestLoadScopes_EmptyPath(t *testing.T) {
	r, err := envscope.LoadScopes("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Names()) != 0 {
		t.Error("expected empty resolver")
	}
}

func TestLoadScopes_MissingFile(t *testing.T) {
	_, err := envscope.LoadScopes("/nonexistent/scopes.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadScopes_InvalidYAML(t *testing.T) {
	p := writeScopeFile(t, ":::invalid yaml:::")
	_, err := envscope.LoadScopes(p)
	if err == nil {
		t.Fatal("expected error for invalid yaml")
	}
}
