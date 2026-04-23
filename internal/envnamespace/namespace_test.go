package envnamespace_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/envnamespace"
)

var baseNamespaces = []envnamespace.Namespace{
	{Name: "prod", Prefix: "PROD_"},
	{Name: "staging", Prefix: "STG_"},
}

func TestNewResolver_Valid(t *testing.T) {
	r, err := envnamespace.NewResolver(baseNamespaces)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}

func TestNewResolver_MissingName(t *testing.T) {
	_, err := envnamespace.NewResolver([]envnamespace.Namespace{
		{Name: "", Prefix: "X_"},
	})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNewResolver_MissingPrefix(t *testing.T) {
	_, err := envnamespace.NewResolver([]envnamespace.Namespace{
		{Name: "dev", Prefix: ""},
	})
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestQualify_AddsPrefix(t *testing.T) {
	r, _ := envnamespace.NewResolver(baseNamespaces)
	input := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	out, err := r.Qualify("prod", input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PROD_DB_HOST"] != "localhost" {
		t.Errorf("expected PROD_DB_HOST=localhost, got %q", out["PROD_DB_HOST"])
	}
	if out["PROD_DB_PORT"] != "5432" {
		t.Errorf("expected PROD_DB_PORT=5432, got %q", out["PROD_DB_PORT"])
	}
}

func TestQualify_UnknownNamespace(t *testing.T) {
	r, _ := envnamespace.NewResolver(baseNamespaces)
	_, err := r.Qualify("unknown", map[string]string{"KEY": "val"})
	if err == nil {
		t.Fatal("expected error for unknown namespace")
	}
}

func TestStrip_RemovesPrefix(t *testing.T) {
	r, _ := envnamespace.NewResolver(baseNamespaces)
	input := map[string]string{"STG_API_KEY": "secret", "STG_TIMEOUT": "30"}
	out, err := r.Strip("staging", input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY=secret, got %q", out["API_KEY"])
	}
	if out["TIMEOUT"] != "30" {
		t.Errorf("expected TIMEOUT=30, got %q", out["TIMEOUT"])
	}
}

func TestStrip_PassesThroughUnprefixed(t *testing.T) {
	r, _ := envnamespace.NewResolver(baseNamespaces)
	input := map[string]string{"OTHER_KEY": "value"}
	out, err := r.Strip("prod", input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["OTHER_KEY"] != "value" {
		t.Errorf("expected OTHER_KEY=value, got %q", out["OTHER_KEY"])
	}
}

func TestNames_ReturnsAll(t *testing.T) {
	r, _ := envnamespace.NewResolver(baseNamespaces)
	names := r.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}
