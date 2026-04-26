# envpriority

The `envpriority` package merges multiple environment variable sources according to
explicit numeric priority levels.

## Concepts

- **Source** — a named set of key/value pairs with an integer priority (1 = highest).
- **Merge** — combines all sources so that for any key the value from the highest-priority source wins.
- **Origin** — the result records which source name contributed each final key.

## Usage

```go
sources := []envpriority.Source{
    {Name: "vault",   Priority: 1, Values: vaultSecrets},
    {Name: "dotenv",  Priority: 2, Values: localEnv},
    {Name: "default", Priority: 3, Values: defaults},
}

result, err := envpriority.Merge(sources)
if err != nil {
    log.Fatal(err)
}

fmt.Println(result.Merged)  // final merged map
fmt.Println(result.Origin)  // key -> winning source name

counts := envpriority.Summary(result)
// map["vault":5 "dotenv":2 "default":1]
```

## Rules

- Priority must be >= 1; an empty source name is rejected.
- When two sources share the same priority number, the one listed first in the slice wins.
- Keys present only in lower-priority sources are still included in the output.
