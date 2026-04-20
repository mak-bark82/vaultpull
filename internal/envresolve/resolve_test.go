package envresolve_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envresolve"
	"github.com/your-org/vaultpull/internal/envtransform"
)

var base = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PASSWORD": "secret",
	"API_KEY":     "abc123",
}

func TestResolve_NoOptions(t *testing.T) {
	r := envresolve.New(envresolve.Options{})
	out, err := r.Resolve(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(base) {
		t.Errorf("expected %d keys, got %d", len(base), len(out))
	}
}

func TestResolve_IncludeFilter(t *testing.T) {
	r := envresolve.New(envresolve.Options{
		IncludeKeys: []string{"DB_HOST", "DB_PASSWORD"},
	})
	out, err := r.Resolve(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["API_KEY"]; ok {
		t.Error("API_KEY should have been excluded")
	}
}

func TestResolve_ExcludeFilter(t *testing.T) {
	r := envresolve.New(envresolve.Options{
		ExcludeKeys: []string{"DB_PASSWORD"},
	})
	out, err := r.Resolve(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been excluded")
	}
}

func TestResolve_TransformUpperCase(t *testing.T) {
	r := envresolve.New(envresolve.Options{
		Transforms: map[string]envtransform.Rule{
			"DB_HOST": {Operation: "uppercase"},
		},
	})
	out, err := r.Resolve(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "LOCALHOST" {
		t.Errorf("expected LOCALHOST, got %q", out["DB_HOST"])
	}
}

func TestResolve_NilInput(t *testing.T) {
	r := envresolve.New(envresolve.Options{})
	out, err := r.Resolve(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map for nil input, got %d keys", len(out))
	}
}

func TestResolve_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"KEY": "value"}
	r := envresolve.New(envresolve.Options{
		ExcludeKeys: []string{"KEY"},
	})
	_, _ = r.Resolve(input)
	if _, ok := input["KEY"]; !ok {
		t.Error("Resolve must not mutate the input map")
	}
}
