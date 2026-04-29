package envtrim_test

import (
	"testing"

	"github.com/wndhydrnt/vaultpull/internal/envtrim"
)

var baseEnv = map[string]string{
	"  KEY_A  ": "  value_a  ",
	"KEY_B":     `"quoted"`,
	"key_c":     "  mixed  ",
}

func TestTrim_TrimKeysAndValues(t *testing.T) {
	opts := envtrim.DefaultOptions()
	out := envtrim.Trim(baseEnv, opts)

	if v, ok := out["KEY_A"]; !ok || v != "value_a" {
		t.Errorf("expected KEY_A=value_a, got %q", v)
	}
}

func TestTrim_StripValueQuotes(t *testing.T) {
	opts := envtrim.DefaultOptions()
	opts.StripValueQuotes = true
	out := envtrim.Trim(baseEnv, opts)

	if v := out["KEY_B"]; v != "quoted" {
		t.Errorf("expected KEY_B=quoted, got %q", v)
	}
}

func TestTrim_NormalizeKeys(t *testing.T) {
	opts := envtrim.DefaultOptions()
	opts.NormalizeKeys = true
	out := envtrim.Trim(baseEnv, opts)

	if _, ok := out["KEY_C"]; !ok {
		t.Error("expected key_c to be normalized to KEY_C")
	}
}

func TestTrim_DoesNotMutateInput(t *testing.T) {
	orig := map[string]string{"  FOO  ": "  bar  "}
	opts := envtrim.DefaultOptions()
	envtrim.Trim(orig, opts)

	if _, ok := orig["  FOO  "]; !ok {
		t.Error("Trim must not mutate the original map")
	}
}

func TestTrim_EmptyKeyDropped(t *testing.T) {
	env := map[string]string{"   ": "value"}
	opts := envtrim.DefaultOptions()
	out := envtrim.Trim(env, opts)

	if len(out) != 0 {
		t.Errorf("expected empty map after trimming blank key, got %v", out)
	}
}

func TestTrimValue_SingleQuotes(t *testing.T) {
	opts := envtrim.Options{StripValueQuotes: true}
	got := envtrim.TrimValue("'hello'", opts)
	if got != "hello" {
		t.Errorf("expected hello, got %q", got)
	}
}

func TestTrimValue_NoQuotesToStrip(t *testing.T) {
	opts := envtrim.Options{StripValueQuotes: true}
	got := envtrim.TrimValue("hello", opts)
	if got != "hello" {
		t.Errorf("expected hello unchanged, got %q", got)
	}
}

func TestTrimValue_ShortString(t *testing.T) {
	opts := envtrim.Options{StripValueQuotes: true}
	got := envtrim.TrimValue("'", opts)
	if got != "'" {
		t.Errorf("expected single quote unchanged, got %q", got)
	}
}
