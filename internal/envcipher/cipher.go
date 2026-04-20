// Package envcipher provides symmetric encryption and decryption
// for secret values before writing them to .env files.
package envcipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// ErrInvalidKey is returned when the key length is not 16, 24, or 32 bytes.
var ErrInvalidKey = errors.New("envcipher: key must be 16, 24, or 32 bytes")

// ErrCiphertextTooShort is returned when the ciphertext is too short to be valid.
var ErrCiphertextTooShort = errors.New("envcipher: ciphertext too short")

// Cipher wraps an AES-GCM block cipher for encrypting and decrypting env values.
type Cipher struct {
	gcm cipher.AEAD
}

// New creates a new Cipher using the provided key.
// Key must be 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
func New(key []byte) (*Cipher, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, ErrInvalidKey
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Cipher{gcm: gcm}, nil
}

// Encrypt encrypts plaintext and returns a base64-encoded ciphertext string.
func (c *Cipher) Encrypt(plaintext string) (string, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := c.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt decodes and decrypts a base64-encoded ciphertext string.
func (c *Cipher) Decrypt(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	ns := c.gcm.NonceSize()
	if len(data) < ns {
		return "", ErrCiphertextTooShort
	}
	nonce, ciphertext := data[:ns], data[ns:]
	plaintext, err := c.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// EncryptMap encrypts all values in the given map and returns a new map.
func (c *Cipher) EncryptMap(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		enc, err := c.Encrypt(v)
		if err != nil {
			return nil, err
		}
		out[k] = enc
	}
	return out, nil
}

// DecryptMap decrypts all values in the given map and returns a new map.
func (c *Cipher) DecryptMap(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		dec, err := c.Decrypt(v)
		if err != nil {
			return nil, err
		}
		out[k] = dec
	}
	return out, nil
}
