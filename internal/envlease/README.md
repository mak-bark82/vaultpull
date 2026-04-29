# envlease

The `envlease` package provides time-bounded leases on secret keys pulled from Vault.

A **lease** grants exclusive access to a named key for a specified duration. If
another process attempts to acquire the same key while a valid lease is held, the
request is rejected until the lease expires or is explicitly released.

## Usage

```go
m, err := envlease.New("/var/run/vaultpull/leases.json", nil)
if err != nil {
    log.Fatal(err)
}

// Acquire a 30-minute lease on DB_PASSWORD for owner "deploy-job"
if err := m.Acquire("DB_PASSWORD", "deploy-job", 30*time.Minute); err != nil {
    log.Fatalf("could not acquire lease: %v", err)
}
defer m.Release("DB_PASSWORD", "deploy-job")
```

## Configuration

Load optional YAML config:

```yaml
path: /var/run/vaultpull/leases.json
default_ttl: 1h
```

```go
cfg, _ := envlease.LoadConfig("lease.yaml")
m, _ := envlease.New(cfg.Path, nil)
```

## Notes

- Leases are persisted as JSON; concurrent processes share state via the file.
- `PurgeExpired` can be called periodically to clean up stale entries.
- Pass an empty path to `New` for an in-memory-only manager (useful in tests).
