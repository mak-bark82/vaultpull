# envwatch

Package `envwatch` provides lightweight file-change detection for local `.env` files.

## Overview

A `Watcher` polls one or more files at a configurable interval. When the SHA-256 checksum of a file changes between polls, a `ChangeEvent` is sent on the `Events` channel.

This is useful for detecting when a `vaultpull sync` has written new values so that dependent processes can be notified or restarted.

## Usage

```go
w := envwatch.New(5 * time.Second)

if err := w.Watch(".env"); err != nil {
    log.Fatal(err)
}

w.Start()
defer w.Stop()

for ev := range w.Events {
    fmt.Printf("file changed: %s (%s → %s)\n", ev.Path, ev.OldHash[:8], ev.NewHash[:8])
}
```

## Types

| Type | Description |
|------|-------------|
| `Watcher` | Polls registered files and emits `ChangeEvent` values |
| `ChangeEvent` | Carries the path and old/new checksums for a changed file |
| `FileState` | Internal snapshot of a file's checksum and mod-time |

## Notes

- Polling is used intentionally to remain dependency-free (no `inotify`/`fsevents`).
- The `Events` channel is buffered (capacity 16) to avoid blocking the poll loop.
- Call `Stop()` to cleanly shut down the background goroutine.
