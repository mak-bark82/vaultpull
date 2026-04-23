# envquota

The `envquota` package enforces size and count limits on environment secret maps.

## Features

- Limit the total number of keys
- Limit the maximum length of any key name
- Limit the maximum length of any value
- Load quota rules from a YAML file
- Returns structured `Violation` results for downstream reporting

## Rule file format (YAML)

```yaml
max_keys: 50
max_key_length: 64
max_val_length: 512
```

Set any field to `0` (or omit it) to apply no limit for that dimension.

## Usage

```go
rule, err := envquota.LoadRule("quota.yaml")
if err != nil {
    log.Fatal(err)
}

result, err := envquota.Check(secrets, rule)
if err != nil {
    fmt.Println(result.Summary())
    for _, v := range result.Violations {
        fmt.Println(" -", v)
    }
}
```
