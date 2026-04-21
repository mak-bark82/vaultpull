# envdrift

The `envdrift` package detects configuration drift between secrets stored in
HashiCorp Vault and the values present in a local `.env` file.

## Concepts

| Status | Meaning |
|--------|---------|
| `StatusMatch` | Key exists in both sources with the same value |
| `StatusChanged` | Key exists in both sources but values differ |
| `StatusMissing` | Key is in Vault but absent from the local file |
| `StatusExtra` | Key is in the local file but not present in Vault |

## Usage

```go
vaultSecrets := map[string]string{"DB_PASS": "new", "API_KEY": "abc"}
localSecrets := map[string]string{"DB_PASS": "old", "API_KEY": "abc"}

report := envdrift.Detect(vaultSecrets, localSecrets)

if report.HasDrift() {
    fmt.Println(report.Summary())
    for _, e := range report.Entries {
        if e.Status != envdrift.StatusMatch {
            fmt.Printf("  %s [%v]\n", e.Key, e.Status)
        }
    }
}
```

## Integration

Combine with `internal/envreader` to load the local file and
`internal/vault/client` to fetch live Vault secrets before calling `Detect`.
