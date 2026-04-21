package envmask

import (
	"strings"
	"testing"
)

func TestApply_MasksMatchingKeys(t *testing.T) {
	m := New(DefaultOptions(), []string{"password", "secret", "token"})
	secrets := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_TOKEN":   "abc123xyz",
		"APP_NAME":    "myapp",
	}
	out := m.Apply(secrets)
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", out["APP_NAME"])
	}
	if out["DB_PASSWORD"] == "supersecret" {
		t.Error("expected DB_PASSWORD to be masked")
	}
	if out["API_TOKEN"] == "abc123xyz" {
		t.Error("expected API_TOKEN to be masked")
	}
}

func TestApply_PartialMask_RevealsLastChars(t *testing.T) {
	opts := Options{Mode: MaskPartial, VisibleChars: 4, MaskChar: '*'}
	m := New(opts, []string{"secret"})
	secrets := map[string]string{"MY_SECRET": "abcdefgh"}
	out := m.Apply(secrets)
	masked := out["MY_SECRET"]
	if !strings.HasSuffix(masked, "efgh") {
		t.Errorf("expected suffix 'efgh', got %q", masked)
	}
	if !strings.HasPrefix(masked, "****") {
		t.Errorf("expected masked prefix, got %q", masked)
	}
}

func TestApply_FullMask_ReplacesAll(t *testing.T) {
	opts := Options{Mode: MaskFull, MaskChar: '#'}
	m := New(opts, []string{"key"})
	secrets := map[string]string{"API_KEY": "hello"}
	out := m.Apply(secrets)
	if out["API_KEY"] != "#####" {
		t.Errorf("expected '#####', got %q", out["API_KEY"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	m := New(DefaultOptions(), []string{"token"})
	secrets := map[string]string{"TOKEN": "original"}
	_ = m.Apply(secrets)
	if secrets["TOKEN"] != "original" {
		t.Error("input map was mutated")
	}
}

func TestApply_EmptyValue_ReturnsEmpty(t *testing.T) {
	m := New(DefaultOptions(), []string{"password"})
	secrets := map[string]string{"PASSWORD": ""}
	out := m.Apply(secrets)
	if out["PASSWORD"] != "" {
		t.Errorf("expected empty string, got %q", out["PASSWORD"])
	}
}

func TestApply_ShortValue_NotPanics(t *testing.T) {
	opts := Options{Mode: MaskPartial, VisibleChars: 10, MaskChar: '*'}
	m := New(opts, []string{"secret"})
	secrets := map[string]string{"SECRET": "abc"}
	out := m.Apply(secrets)
	// value shorter than VisibleChars: returned as-is
	if out["SECRET"] != "abc" {
		t.Errorf("expected 'abc', got %q", out["SECRET"])
	}
}
