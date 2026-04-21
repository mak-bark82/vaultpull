# envsnap

The `envsnap` package provides point-in-time snapshots of environment variable maps, with support for persistence, diffing, and equality checks.

## Features

- **Take** — Captures an immutable snapshot of an env map with a timestamp, source label, and SHA-256 checksum.
- **Equal** — Compares two snapshots by checksum for fast equality testing.
- **Diff** — Returns added, removed, and changed keys between two snapshots.
- **Store** — Persists snapshots as JSON files in a directory and retrieves the latest or full history.

## Usage

```go
store, err := envsnap.NewStore(".vaultpull/snaps")
if err != nil { ... }

snap := envsnap.Take("prod/app", secrets, time.Now)
_ = store.Save(snap)

prev, _ := store.Latest()
added, removed, changed := envsnap.Diff(prev, snap)
fmt.Println(snap.Summary())
```

## Snapshot files

Snapshots are stored as `<timestamp>_<source>.snap.json` files (e.g. `20240115T120000Z_prod_app.snap.json`) and are sorted lexicographically to determine recency.

## Notes

- Snapshots are immutable once taken; mutating the original map after `Take` does not affect the snapshot.
- Checksums are order-independent (keys are sorted before hashing).
