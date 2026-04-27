package envencrypt

import (
	"errors"
	"strings"
	"testing"
)

// stubCipher is a trivial reversible cipher for testing.
type stubCipher struct{ failOn string }

func (s *stubCipher) Encrypt(v string) (string, error) {
	if s.failOn == v {
		return "", errors.New("encrypt error")
	}
	return "enc:" + v, nil
}

func (s *stubCipher) Decrypt(v string) (string, error) {
	if s.failOn == v {
		return "", errors.New("decrypt error")
	}
	if !strings.HasPrefix(v, "enc:") {
		return "", errors.New("not encrypted")
	}
	return strings.TrimPrefix(v, "enc:"), nil
}

func TestNew_NilCipher(t *testing.T) {
	_, err := New(nil, []string{"SECRET"})
	if err == nil {
		t.Fatal("expected error for nil cipher")
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New(&stubCipher{}, []string{""})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidRegexp(t *testing.T) {
	_, err := New(&stubCipher{}, []string{"[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regexp")
	}
}

func TestEncrypt_MatchingKeys(t *testing.T) {
	p, err := New(&stubCipher{}, []string{"SECRET", "TOKEN"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	env := map[string]string{
		"DB_SECRET": "mysecret",
		"API_TOKEN": "tok123",
		"HOST":      "localhost",
	}
	out, err := p.Encrypt(env)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if out["DB_SECRET"] != "enc:mysecret" {
		t.Errorf("DB_SECRET: got %q", out["DB_SECRET"])
	}
	if out["API_TOKEN"] != "enc:tok123" {
		t.Errorf("API_TOKEN: got %q", out["API_TOKEN"])
	}
	if out["HOST"] != "localhost" {
		t.Errorf("HOST should be unchanged, got %q", out["HOST"])
	}
}

func TestDecrypt_Roundtrip(t *testing.T) {
	p, _ := New(&stubCipher{}, []string{"SECRET"})
	env := map[string]string{"DB_SECRET": "mysecret", "HOST": "localhost"}
	enc, err := p.Encrypt(env)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	dec, err := p.Decrypt(enc)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if dec["DB_SECRET"] != "mysecret" {
		t.Errorf("roundtrip failed: got %q", dec["DB_SECRET"])
	}
}

func TestEncrypt_CipherError_PropagatesKey(t *testing.T) {
	p, _ := New(&stubCipher{failOn: "bad"}, []string{"KEY"})
	_, err := p.Encrypt(map[string]string{"MY_KEY": "bad"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "MY_KEY") {
		t.Errorf("error should mention key, got: %v", err)
	}
}

func TestEncrypt_DoesNotMutateInput(t *testing.T) {
	p, _ := New(&stubCipher{}, []string{"SECRET"})
	env := map[string]string{"DB_SECRET": "original"}
	p.Encrypt(env)
	if env["DB_SECRET"] != "original" {
		t.Error("input map was mutated")
	}
}
