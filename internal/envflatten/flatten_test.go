package envflatten_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/envflatten"
)

func TestFlatten_SimpleKeys(t *testing.T) {
	input := map[string]interface{}{
		"host": "localhost",
		"port": "5432",
	}
	out, err := envflatten.Flatten(input, envflatten.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", out["HOST"])
	}
	if out["PORT"] != "5432" {
		t.Errorf("expected PORT=5432, got %q", out["PORT"])
	}
}

func TestFlatten_NestedMap(t *testing.T) {
	input := map[string]interface{}{
		"db": map[string]interface{}{
			"host": "db.local",
			"port": "5432",
		},
	}
	out, err := envflatten.Flatten(input, envflatten.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "db.local" {
		t.Errorf("expected DB_HOST=db.local, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", out["DB_PORT"])
	}
}

func TestFlatten_DeeplyNested(t *testing.T) {
	input := map[string]interface{}{
		"app": map[string]interface{}{
			"cache": map[string]interface{}{
				"ttl": "300",
			},
		},
	}
	out, err := envflatten.Flatten(input, envflatten.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_CACHE_TTL"] != "300" {
		t.Errorf("expected APP_CACHE_TTL=300, got %q", out["APP_CACHE_TTL"])
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	input := map[string]interface{}{"name": "vault"}
	opts := envflatten.Options{Separator: "_", UpperCase: true, Prefix: "APP"}
	out, err := envflatten.Flatten(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_NAME"] != "vault" {
		t.Errorf("expected APP_NAME=vault, got %q", out["APP_NAME"])
	}
}

func TestFlatten_NilValue(t *testing.T) {
	input := map[string]interface{}{"key": nil}
	out, err := envflatten.Flatten(input, envflatten.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "" {
		t.Errorf("expected KEY='', got %q", out["KEY"])
	}
}

func TestFlatten_LowercaseOption(t *testing.T) {
	input := map[string]interface{}{"Host": "localhost"}
	opts := envflatten.Options{Separator: "_", UpperCase: false}
	out, err := envflatten.Flatten(input, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["Host"] != "localhost" {
		t.Errorf("expected Host=localhost, got %q", out["Host"])
	}
}
