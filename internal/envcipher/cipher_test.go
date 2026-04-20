package envcipher_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/envcipher"
)

var validKey = []byte("0123456789abcdef") // 16 bytes

func TestNew_ValidKey(t *testing.T) {
	_, err := envcipher.New(validKey)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNew_InvalidKey(t *testing.T) {
	_, err := envcipher.New([]byte("short"))
	if err == nil {
		t.Fatal("expected error for invalid key length")
	}
}

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	c, _ := envcipher.New(validKey)
	plaintext := "super-secret-value"

	enc, err := c.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}
	if enc == plaintext {
		t.Fatal("encrypted value should differ from plaintext")
	}

	dec, err := c.Decrypt(enc)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	if dec != plaintext {
		t.Errorf("expected %q, got %q", plaintext, dec)
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	c, _ := envcipher.New(validKey)
	_, err := c.Decrypt("!!!not-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestDecrypt_TooShort(t *testing.T) {
	c, _ := envcipher.New(validKey)
	// base64 of a very short byte slice
	_, err := c.Decrypt("YWJj")
	if err == nil {
		t.Fatal("expected error for ciphertext too short")
	}
}

func TestEncryptMap_DecryptMap_Roundtrip(t *testing.T) {
	c, _ := envcipher.New(validKey)
	secrets := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abc123",
		"TOKEN":       "tok_live_xyz",
	}

	encrypted, err := c.EncryptMap(secrets)
	if err != nil {
		t.Fatalf("EncryptMap failed: %v", err)
	}
	for k, v := range encrypted {
		if v == secrets[k] {
			t.Errorf("key %s: encrypted value should not equal plaintext", k)
		}
	}

	decrypted, err := c.DecryptMap(encrypted)
	if err != nil {
		t.Fatalf("DecryptMap failed: %v", err)
	}
	for k, want := range secrets {
		if got := decrypted[k]; got != want {
			t.Errorf("key %s: expected %q, got %q", k, want, got)
		}
	}
}

func TestEncrypt_ProducesUniqueOutputs(t *testing.T) {
	c, _ := envcipher.New(validKey)
	a, _ := c.Encrypt("same-value")
	b, _ := c.Encrypt("same-value")
	// GCM uses a random nonce so each call should differ
	if strings.EqualFold(a, b) {
		t.Error("expected different ciphertexts for same plaintext due to random nonce")
	}
}
