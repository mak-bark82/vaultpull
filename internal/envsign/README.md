# envSign

Package `envSign` provides HMAC-SHA256 signing and verification for env secret maps, plus persistent storage of signature records.

## Overview

This module protects against silent tampering of `.env` files or in-memory secret maps by computing a deterministic, keyed signature over all key-value pairs.

## Usage

### Signing

```go
signer, err := envSign.New([]byte(os.Getenv("SIGN_KEY")))
if err != nil { ... }

sig := signer.Sign(secrets)
```

### Verification

```go
if err := signer.Verify(secrets, sig); err != nil {
    log.Fatal("secrets have been tampered with:", err)
}
```

### Persisting a Record

```go
rec := envSign.Record{
    File:      ".env.production",
    Signature: sig,
    SignedAt:  time.Now(),
}
envSign.SaveRecord(".vault/sig.json", rec)
```

### Loading a Record

```go
rec, err := envSign.LoadRecord(".vault/sig.json")
```

## Notes

- Signatures are order-independent; the canonical form sorts keys before hashing.
- Records are stored as JSON with `0600` permissions.
- Use a strong, randomly generated key stored in Vault or a secure keyring.
