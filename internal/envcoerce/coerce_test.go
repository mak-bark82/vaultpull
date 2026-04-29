package envcoerce_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envcoerce"
)

var baseEnv = map[string]string{
	"ENABLED":  "TRUE",
	"PORT":     " 8080 ",
	"RATIO":    "3.14",
	"APP_NAME": "  myapp  ",
	"UNKNOWN":  "whatever",
}

func TestCoerce_Bool(t *testing.T) {
	rules := []envcoerce.Rule{{Key: "ENABLED", Type: envcoerce.TypeBool}}
	out, results := envcoerce.Coerce(baseEnv, rules)
	if out["ENABLED"] != "true" {
		t.Errorf("expected \"true\", got %q", out["ENABLED"])
	}
	if len(results) != 1 || results[0].Err != nil {
		t.Errorf("unexpected result: %+v", results)
	}
	if !results[0].Changed {
		t.Error("expected Changed=true for bool normalisation")
	}
}

func TestCoerce_Int(t *testing.T) {
	rules := []envcoerce.Rule{{Key: "PORT", Type: envcoerce.TypeInt}}
	out, results := envcoerce.Coerce(baseEnv, rules)
	if out["PORT"] != "8080" {
		t.Errorf("expected \"8080\", got %q", out["PORT"])
	}
	if results[0].Err != nil {
		t.Errorf("unexpected error: %v", results[0].Err)
	}
}

func TestCoerce_Float(t *testing.T) {
	rules := []envcoerce.Rule{{Key: "RATIO", Type: envcoerce.TypeFloat}}
	out, results := envcoerce.Coerce(baseEnv, rules)
	if out["RATIO"] != "3.14" {
		t.Errorf("expected \"3.14\", got %q", out["RATIO"])
	}
	if results[0].Err != nil {
		t.Errorf("unexpected error: %v", results[0].Err)
	}
}

func TestCoerce_String_TrimsSpace(t *testing.T) {
	rules := []envcoerce.Rule{{Key: "APP_NAME", Type: envcoerce.TypeString}}
	out, results := envcoerce.Coerce(baseEnv, rules)
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected \"myapp\", got %q", out["APP_NAME"])
	}
	if !results[0].Changed {
		t.Error("expected Changed=true after trimming")
	}
}

func TestCoerce_InvalidBool_ReturnsError(t *testing.T) {
	env := map[string]string{"FLAG": "notabool"}
	rules := []envcoerce.Rule{{Key: "FLAG", Type: envcoerce.TypeBool}}
	out, results := envcoerce.Coerce(env, rules)
	if results[0].Err == nil {
		t.Error("expected error for invalid bool")
	}
	// original value should be preserved on error
	if out["FLAG"] != "notabool" {
		t.Errorf("expected original value preserved, got %q", out["FLAG"])
	}
}

func TestCoerce_MissingKey_Skipped(t *testing.T) {
	rules := []envcoerce.Rule{{Key: "DOES_NOT_EXIST", Type: envcoerce.TypeInt}}
	_, results := envcoerce.Coerce(baseEnv, rules)
	if len(results) != 0 {
		t.Errorf("expected no results for missing key, got %d", len(results))
	}
}

func TestCoerce_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"PORT": " 9090 "}
	rules := []envcoerce.Rule{{Key: "PORT", Type: envcoerce.TypeInt}}
	envcoerce.Coerce(env, rules)
	if env["PORT"] != " 9090 " {
		t.Error("input map was mutated")
	}
}

func TestCoerce_UnknownType_ReturnsError(t *testing.T) {
	env := map[string]string{"X": "val"}
	rules := []envcoerce.Rule{{Key: "X", Type: envcoerce.TargetType("hex")}}
	_, results := envcoerce.Coerce(env, rules)
	if results[0].Err == nil {
		t.Error("expected error for unknown type")
	}
}
