# envaudit

The `envaudit` package provides structured audit recording for vaultpull sync operations.

## Overview

Every key-level change during a sync run (added, changed, removed, skipped) can be recorded as an `Event`. The `Recorder` collects these events in memory and can export them in multiple formats.

## Usage

```go
r := envaudit.New(nil) // nil uses time.Now

r.Record(envaudit.EventAdded,   "DB_HOST",  ".env",      "")
r.Record(envaudit.EventChanged, "API_KEY",  ".env.prod", "rotated")
r.Record(envaudit.EventSkipped, "OLD_FLAG", ".env",      "no-overwrite")

// Human-readable summary
fmt.Println(r.Summary())

// JSON export
r.ExportJSON(os.Stdout)

// CSV export
r.ExportCSV(os.Stdout)
```

## Event Kinds

| Kind      | Meaning                                      |
|-----------|----------------------------------------------|
| `added`   | Key did not exist locally; written from Vault |
| `changed` | Key existed locally; value updated from Vault |
| `removed` | Key removed during rotation or promotion      |
| `skipped` | Key intentionally left unchanged              |

## Notes

- `Events()` always returns a copy; callers cannot mutate internal state.
- A custom clock can be injected for deterministic testing.
