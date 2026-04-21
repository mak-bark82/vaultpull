# envttl

The `envttl` package provides time-to-live (TTL) tracking for individual secret keys managed by vaultpull.

## Overview

When secrets are pulled from Vault, each key can be assigned an expiry duration. The `Store` persists TTL records to a JSON file and lets callers check whether a key has elapsed, list all expired keys, or remove records.

## Usage

```go
store, err := envttl.New(".vaultpull/ttl.json")
if err != nil {
    log.Fatal(err)
}

// Register a key with a 24-hour TTL
if err := store.Set("DB_PASSWORD", 24*time.Hour); err != nil {
    log.Fatal(err)
}

// Check expiry before using a cached value
if store.IsExpired("DB_PASSWORD") {
    // re-pull from Vault
}

// List all expired keys
for _, key := range store.Expired() {
    fmt.Println("expired:", key)
}

// Remove a record when a key is deleted
_ = store.Remove("DB_PASSWORD")
```

## File Format

TTL records are stored as a JSON object keyed by secret name:

```json
{
  "DB_PASSWORD": {
    "key": "DB_PASSWORD",
    "expires_at": "2024-06-01T12:00:00Z"
  }
}
```

The file is written with `0600` permissions and its parent directory is created automatically.
