# envclone

The `envclone` package provides functionality to clone (copy) secrets from one
environment map into another, with fine-grained control over conflict resolution
and dry-run preview.

## Usage

```go
src := map[string]string{"DB_HOST": "prod-db", "API_KEY": "abc123"}
dst := map[string]string{"DB_HOST": "staging-db"}

result, err := envclone.Clone(src, dst, envclone.CloneOptions{
    Overwrite: false,
    DryRun:    false,
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(result.Summary())
// copied=1 skipped=1 overwritten=0
```

## Options

| Option      | Type   | Description                                          |
|-------------|--------|------------------------------------------------------|
| `Overwrite` | `bool` | Replace existing keys in dst with values from src.   |
| `DryRun`    | `bool` | Report changes without modifying dst.                |

## Result Fields

- `Copied` — keys added to dst (new keys from src).
- `Skipped` — keys present in dst that were not overwritten.
- `Overwrite` — keys in dst that were replaced by src values.

Use `Result.Summary()` for a concise one-line status string.
