package envscope_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envscope"
)

func baseScopes() []envscope.Scope {
	return []envscope.Scope{
		{Name: "dev", Prefix: "secret/dev"},
		{Name: "prod", Prefix: "secret/prod"},
	}
}

func TestNewResolver_Valid(t *testing.T) {
	_, err := envscope.NewResolver(baseScopes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewResolver_MissingName(t *testing.T) {
	_, err := envscope.NewResolver([]envscope.Scope{{Prefix: "secret/x"}})
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestNewResolver_MissingPrefix(t *testing.T) {
	_, err := envscope.NewResolver([]envscope.Scope{{Name: "dev"}})
	if err == nil {
		t.Fatal("expected error for missing prefix")
	}
}

func TestResolve_KnownScope(t *testing.T) {
	r, _ := envscope.NewResolver(baseScopes())
	got, err := r.Resolve("dev", "myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "secret/dev/myapp/config"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolve_UnknownScope(t *testing.T) {
	r, _ := envscope.NewResolver(baseScopes())
	_, err := r.Resolve("staging", "myapp/config")
	if err == nil {
		t.Fatal("expected error for unknown scope")
	}
}

func TestNames_ReturnsAll(t *testing.T) {
	r, _ := envscope.NewResolver(baseScopes())
	names := r.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}
