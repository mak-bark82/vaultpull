package envtransform_test

import (
	"testing"

	"github.com/nicholasgasior/vaultpull/internal/envtransform"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_PASS": "  secret  ",
		"API_KEY": "abc123",
	}
}

func TestTransform_UpperCase(t *testing.T) {
	out := envtransform.Transform(baseSecrets(), envtransform.Rule{UpperCase: true})
	if out["API_KEY"] != "ABC123" {
		t.Errorf("expected ABC123, got %s", out["API_KEY"])
	}
}

func TestTransform_LowerCase(t *testing.T) {
	out := envtransform.Transform(map[string]string{"X": "HELLO"}, envtransform.Rule{LowerCase: true})
	if out["X"] != "hello" {
		t.Errorf("expected hello, got %s", out["X"])
	}
}

func TestTransform_TrimSpace(t *testing.T) {
	out := envtransform.Transform(baseSecrets(), envtransform.Rule{TrimSpace: true})
	if out["DB_PASS"] != "secret" {
		t.Errorf("expected 'secret', got %q", out["DB_PASS"])
	}
}

func TestTransform_PrefixSuffix(t *testing.T) {
	out := envtransform.Transform(map[string]string{"K": "val"}, envtransform.Rule{Prefix: "pre_", Suffix: "_suf"})
	if out["K"] != "pre_val_suf" {
		t.Errorf("expected pre_val_suf, got %s", out["K"])
	}
}

func TestTransform_NoRules(t *testing.T) {
	in := map[string]string{"A": "unchanged"}
	out := envtransform.Transform(in, envtransform.Rule{})
	if out["A"] != "unchanged" {
		t.Errorf("expected unchanged, got %s", out["A"])
	}
}

func TestTransformValue_Direct(t *testing.T) {
	v := envtransform.TransformValue("  Hello  ", envtransform.Rule{TrimSpace: true, UpperCase: true})
	if v != "HELLO" {
		t.Errorf("expected HELLO, got %s", v)
	}
}

func TestTransform_DoesNotMutateInput(t *testing.T) {
	in := map[string]string{"K": "original"}
	envtransform.Transform(in, envtransform.Rule{Prefix: "x"})
	if in["K"] != "original" {
		t.Error("input map was mutated")
	}
}
