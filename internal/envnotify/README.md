# envnotify

Provides structured notification events for secret sync operations.

## Overview

`envnotify` exposes two types:

- **`Notifier`** — collects and writes `Event` records (INFO / WARN / ERROR) to any `io.Writer`.
- **`Dispatcher`** — bridges `diff.Change` slices to `Notifier` events, mapping change types to appropriate severity levels.

## Usage

```go
import (
    "os"
    "github.com/dstotijn/vaultpull/internal/envnotify"
    "github.com/dstotijn/vaultpull/internal/diff"
)

n := envnotify.New(os.Stdout)
d := envnotify.NewDispatcher(n)

changes := diff.Compare(existing, incoming)
d.Dispatch(changes)

fmt.Println(n.Summary())
```

## Levels

| Change Type | Level  |
|-------------|--------|
| Added       | INFO   |
| Changed     | WARN   |
| Removed     | WARN   |
| Unchanged   | INFO   |

## Notes

- Passing `nil` as the writer to `New` disables output but still records events.
- `Events()` always returns a copy to prevent external mutation.
