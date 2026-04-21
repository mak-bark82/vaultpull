package envSign_test

import (
	"testing"

	envSign "github.com/yourusername/vaultpull/internal/envSign"
)

func TestNew_ValidKey(t *testing.T) {
	s, err := envSign.New([]byte("secret"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil signer")
	}
}

func TestNew_EmptyKey(t *testing.T) {
	_, err := envSign.New([]byte{})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestSign_IsStable(t *testing.T) {
	s, _ := envSign.New([]byte("mykey"))
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}

	sig1 := s.Sign(secrets)
	sig2 := s.Sign(secrets)
	if sig1 != sig2 {
		t.Errorf("expected stable signature, got %q and %q", sig1, sig2)
	}
}

func TestSign_OrderIndependent(t *testing.T) {
	s, _ := envSign.New([]byte("mykey"))

	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAR": "2", "FOO": "1"}

	if s.Sign(a) != s.Sign(b) {
		t.Error("expected same signature regardless of map iteration order")
	}
}

func TestVerify_Valid(t *testing.T) {
	s, _ := envSign.New([]byte("mykey"))
	secrets := map[string]string{"DB_PASS": "secret123"}

	sig := s.Sign(secrets)
	if err := s.Verify(secrets, sig); err != nil {
		t.Errorf("expected valid verification, got %v", err)
	}
}

func TestVerify_TamperedValue(t *testing.T) {
	s, _ := envSign.New([]byte("mykey"))
	original := map[string]string{"DB_PASS": "original"}
	sig := s.Sign(original)

	tampered := map[string]string{"DB_PASS": "hacked"}
	if err := s.Verify(tampered, sig); err == nil {
		t.Error("expected verification to fail for tampered value")
	}
}

func TestVerify_TamperedKey(t *testing.T) {
	s, _ := envSign.New([]byte("mykey"))
	original := map[string]string{"FOO": "bar"}
	sig := s.Sign(original)

	tampered := map[string]string{"FOO": "bar", "EXTRA": "injected"}
	if err := s.Verify(tampered, sig); err == nil {
		t.Error("expected verification to fail when key is injected")
	}
}

func TestVerify_WrongKey(t *testing.T) {
	signer1, _ := envSign.New([]byte("key-one"))
	signer2, _ := envSign.New([]byte("key-two"))

	secrets := map[string]string{"TOKEN": "abc"}
	sig := signer1.Sign(secrets)

	if err := signer2.Verify(secrets, sig); err == nil {
		t.Error("expected verification to fail with different key")
	}
}
