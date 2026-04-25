# envretry

Provides configurable retry logic for transient vault operations.

## Usage

```go
import "github.com/user/vaultpull/internal/envretry"

// Use the default policy (3 attempts, 200ms base delay, 2x backoff)
p := envretry.DefaultPolicy()

result := envretry.Do(p, func() (retryable bool, err error) {
    err = callVault()
    if isTransient(err) {
        return true, err   // will be retried
    }
    return false, err      // fatal — stop immediately
})

if result.Err != nil {
    log.Fatalf("vault call failed after %d attempts: %v", result.Attempts, result.Err)
}
```

## Policy file (YAML)

```yaml
max_attempts: 5
delay_ms: 300
multiplier: 2.0
```

Load with:

```go
p, err := envretry.LoadPolicy("retry.yaml")
```

If the path is empty, `DefaultPolicy` is returned.

## Fields

| Field | Description |
|---|---|
| `max_attempts` | Total number of attempts (>= 1) |
| `delay_ms` | Initial delay between attempts in milliseconds |
| `multiplier` | Backoff multiplier applied after each failure |
