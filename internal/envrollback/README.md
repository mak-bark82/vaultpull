# envrollback

The `envrollback` package provides snapshot-based rollback for `.env` files managed by `vaultpull`.

## Overview

Before overwriting a local `.env` file, callers can save a timestamped snapshot. If the new secrets cause issues, the previous state can be restored in one call.

## Types

### `Rollbacker`

Created with `New(snapshotDir string)`. Snapshots are stored as plain-text files inside `snapshotDir`.

### `Snapshot`

Holds the file path, timestamp, and key/value map of a saved state.

### `RestoreResult`

Returned by `Restore`, reporting how many keys were written back.

## Usage

```go
r, err := envrollback.New(".vaultpull/snapshots")
if err != nil { ... }

// Save current state before syncing
_ = r.Save(".env", currentEnv, time.Now())

// Later, roll back to the last known-good state
snap, err := r.Latest(".env")
if err != nil { ... }
result, err := envrollback.Restore(".env", snap)
fmt.Printf("Restored %d keys to %s\n", result.Written, result.File)
```

## Notes

- Snapshot files are stored with mode `0600`.
- `Latest` returns `nil, nil` when no snapshots exist for the target file.
- Keys are written in sorted order for deterministic diffs.
