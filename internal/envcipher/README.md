# envcipher

Provides AES-GCM symmetric encryption and decryption for secret values managed by `vaultpull`.

## Usage

```go
import "github.com/yourusername/vaultpull/internal/envcipher"

// Key must be 16, 24, or 32 bytes (AES-128 / 192 / 256)
key := []byte("0123456789abcdef")

c, err := envcipher.New(key)
if err != nil {
    log.Fatal(err)
}

// Encrypt a single value
enc, err := c.Encrypt("my-secret")

// Decrypt it back
plain, err := c.Decrypt(enc)

// Encrypt / decrypt an entire secrets map
encMap, err := c.EncryptMap(secrets)
decMap, err := c.DecryptMap(encMap)
```

## Notes

- Each call to `Encrypt` generates a fresh random nonce, so identical plaintexts produce different ciphertexts.
- Ciphertexts are base64-encoded for safe storage in text-based `.env` files.
- The key is never stored or logged by this package; callers are responsible for secure key management (e.g. sourcing from Vault itself or a KMS).
