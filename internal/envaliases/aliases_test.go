package envaliases_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/envaliases"
)

func baseAliases() []envaliases.Alias {
	return []envaliases.Alias{
		{Name: "db", Keys: []string{"DB_HOST", "DB_PORT", "DB_NAME"}},
		{Name: "cache", Keys: []string{"REDIS_HOST", "REDIS_PORT"}},
	}
}

func TestNewResolver_Valid(t *testing.T) {
	_, err := envaliases.NewResolver(baseAliases())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewResolver_EmptyName(t *testing.T) {
	_, err := envaliases.NewResolver([]envaliases.Alias{{Name: "", Keys: []string{"FOO"}}})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNewResolver_NoKeys(t *testing.T) {
	_, err := envaliases.NewResolver([]envaliases.Alias{{Name: "empty", Keys: nil}})
	if err == nil {
		t.Fatal("expected error for empty keys")
	}
}

func TestResolve_KnownAlias(t *testing.T) {
	r, _ := envaliases.NewResolver(baseAliases())
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	out, err := r.Resolve("db", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" || out["DB_PORT"] != "5432" {
		t.Errorf("unexpected output: %v", out)
	}
	if _, ok := out["DB_NAME"]; ok {
		t.Error("DB_NAME should be skipped when missing from env")
	}
}

func TestResolve_UnknownAlias(t *testing.T) {
	r, _ := envaliases.NewResolver(baseAliases())
	_, err := r.Resolve("nope", map[string]string{})
	if err == nil {
		t.Fatal("expected error for unknown alias")
	}
}

func TestExpand_MixedKeys(t *testing.T) {
	r, _ := envaliases.NewResolver(baseAliases())
	keys := r.Expand([]string{"db", "APP_KEY"})
	if len(keys) != 4 {
		t.Fatalf("expected 4 keys, got %d: %v", len(keys), keys)
	}
}

func TestLoadAliases_EmptyPath(t *testing.T) {
	_, err := envaliases.LoadAliases("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadAliases_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "aliases.yaml")
	content := "aliases:\n  - name: svc\n    keys:\n      - SVC_HOST\n      - SVC_PORT\n"
	_ = os.WriteFile(path, []byte(content), 0644)
	r, err := envaliases.LoadAliases(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := r.Resolve("svc", map[string]string{"SVC_HOST": "api.local", "SVC_PORT": "8080"})
	if err != nil || len(out) != 2 {
		t.Errorf("unexpected result: %v %v", out, err)
	}
}
