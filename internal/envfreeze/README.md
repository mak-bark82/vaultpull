# envfreeze

`envfreeze` provides a mechanism to mark specific environment variable keys as **frozen** — protected against accidental modification during sync or merge operations.

Frozen keys are persisted to a JSON file and can be checked before applying changes.

## Usage

```go
f, err := envfreeze.New(".vaultpull/freeze.json")
if err != nil {
    log.Fatal(err)
}

// Freeze keys
_ = f.Freeze([]string{"DB_PASSWORD", "API_SECRET"}, "locked by ops", nil)

// Check before writing
if f.IsFrozen("DB_PASSWORD") {
    fmt.Println("DB_PASSWORD is frozen — skipping update")
}

// Remove a freeze
_ = f.Unfreeze("API_SECRET")
```

## Freeze file format

```json
{
  "keys": ["API_SECRET", "DB_PASSWORD"],
  "frozen_at": "2024-01-01T00:00:00Z",
  "comment": "locked by ops"
}
```

## Notes

- Keys are stored sorted for deterministic diffs.
- An empty or missing freeze file is treated as no frozen keys.
- Empty key strings are silently ignored.
