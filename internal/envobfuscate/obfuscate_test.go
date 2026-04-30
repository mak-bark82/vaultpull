package envobfuscate_test

import (
	"testing"

	"github.com/nicholasgasior/vaultpull/internal/envobfuscate"
)

var baseEnv = map[string]string{
	"DB_PASSWORD":  "supersecret",
	"API_KEY":      "abc123",
	"APP_NAME":     "vaultpull",
	"SECRET_TOKEN": "tok-xyz",
}

func TestNew_Valid(t *testing.T) {
	_, err := envobfuscate.New([]envobfuscate.Rule{
		{Pattern: "PASSWORD", Strategy: "mask"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	_, err := envobfuscate.New([]envobfuscate.Rule{
		{Pattern: "", Strategy: "mask"},
	})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_NoRules(t *testing.T) {
	_, err := envobfuscate.New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := envobfuscate.New([]envobfuscate.Rule{
		{Pattern: "[invalid", Strategy: "mask"},
	})
	if err == nil {
		t.Fatal("expected error for invalid regexp")
	}
}

func TestApply_MaskStrategy(t *testing.T) {
	o, _ := envobfuscate.New([]envobfuscate.Rule{
		{Pattern: "PASSWORD", Strategy: "mask", RevealChars: 3},
	})
	res := o.Apply(baseEnv)
	got := res.Env["DB_PASSWORD"]
	if len(got) != len("supersecret") {
		t.Errorf("expected same length, got %q", got)
	}
	if got[len(got)-3:] != "ret" {
		t.Errorf("expected last 3 chars revealed, got %q", got)
	}
}

func TestApply_RemoveStrategy(t *testing.T) {
	o, _ := envobfuscate.New([]envobfuscate.Rule{
		{Pattern: "SECRET", Strategy: "remove"},
	})
	res := o.Apply(baseEnv)
	if res.Env["SECRET_TOKEN"] != "" {
		t.Errorf("expected empty value, got %q", res.Env["SECRET_TOKEN"])
	}
}

func TestApply_HashStrategy(t *testing.T) {
	o, _ := envobfuscate.New([]envobfuscate.Rule{
		{Pattern: "API_KEY", Strategy: "hash"},
	})
	res := o.Apply(baseEnv)
	got := res.Env["API_KEY"]
	if got == "abc123" {
		t.Error("expected value to be obfuscated")
	}
	if got != "[redacted:6]" {
		t.Errorf("unexpected hash value: %q", got)
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	o, _ := envobfuscate.New([]envobfuscate.Rule{
		{Pattern: "PASSWORD", Strategy: "mask"},
	})
	original := "supersecret"
	env := map[string]string{"DB_PASSWORD": original}
	o.Apply(env)
	if env["DB_PASSWORD"] != original {
		t.Error("input map was mutated")
	}
}

func TestApply_UnmatchedKeyUnchanged(t *testing.T) {
	o, _ := envobfuscate.New([]envobfuscate.Rule{
		{Pattern: "SECRET", Strategy: "mask"},
	})
	res := o.Apply(baseEnv)
	if res.Env["APP_NAME"] != "vaultpull" {
		t.Errorf("expected APP_NAME unchanged, got %q", res.Env["APP_NAME"])
	}
}

func TestApply_ChangedListPopulated(t *testing.T) {
	o, _ := envobfuscate.New([]envobfuscate.Rule{
		{Pattern: "PASSWORD|API_KEY", Strategy: "remove"},
	})
	res := o.Apply(baseEnv)
	if len(res.Changed) != 2 {
		t.Errorf("expected 2 changed keys, got %d", len(res.Changed))
	}
}
