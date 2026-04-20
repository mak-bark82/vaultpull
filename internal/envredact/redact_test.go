package envredact_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envredact"
)

func TestIsSensitive_MatchesDefaultPatterns(t *testing.T) {
	r := envredact.New(nil, "")
	cases := []struct {
		key       string
		expected  bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"AUTH_TOKEN", true},
		{"APP_NAME", false},
		{"PORT", false},
	}
	for _, tc := range cases {
		got := r.IsSensitive(tc.key)
		if got != tc.expected {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.expected)
		}
	}
}

func TestRedact_MasksSensitiveValues(t *testing.T) {
	r := envredact.New(nil, "REDACTED")
	secrets := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"APP_NAME":    "myapp",
		"API_KEY":     "abc123",
	}
	out := r.Redact(secrets)
	if out["DB_PASSWORD"] != "REDACTED" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", out["DB_PASSWORD"])
	}
	if out["API_KEY"] != "REDACTED" {
		t.Errorf("expected API_KEY to be redacted, got %q", out["API_KEY"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", out["APP_NAME"])
	}
}

func TestRedact_DoesNotMutateInput(t *testing.T) {
	r := envredact.New(nil, "")
	orig := map[string]string{"SECRET_KEY": "real-value"}
	_ = r.Redact(orig)
	if orig["SECRET_KEY"] != "real-value" {
		t.Error("Redact mutated the input map")
	}
}

func TestRedact_CustomPatterns(t *testing.T) {
	r := envredact.New([]string{"internal"}, "[hidden]")
	secrets := map[string]string{
		"INTERNAL_HOST": "10.0.0.1",
		"PUBLIC_URL":    "https://example.com",
	}
	out := r.Redact(secrets)
	if out["INTERNAL_HOST"] != "[hidden]" {
		t.Errorf("expected INTERNAL_HOST redacted, got %q", out["INTERNAL_HOST"])
	}
	if out["PUBLIC_URL"] != "https://example.com" {
		t.Errorf("expected PUBLIC_URL unchanged, got %q", out["PUBLIC_URL"])
	}
}
