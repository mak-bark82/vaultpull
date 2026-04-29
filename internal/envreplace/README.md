# envreplace

The `envreplace` package applies regex-based value replacement rules to env maps.

## Usage

```go
rules := []envreplace.Rule{
    {Pattern: "localhost", With: "db.prod.internal"},
    {Key: "API_URL", Pattern: `http://old\.example\.com`, With: "https://api.example.com"},
}

out, results, err := envreplace.Replace(env, rules)
if err != nil {
    log.Fatal(err)
}
fmt.Println(envreplace.Summary(results))
```

## Rule Fields

| Field     | Description                                              |
|-----------|----------------------------------------------------------|
| `key`     | Optional. Restrict replacement to this specific key.     |
| `pattern` | Required. Regular expression to match in the value.      |
| `with`    | Replacement string. Supports `$1` capture group syntax.  |

## Loading from YAML

```yaml
rules:
  - pattern: "localhost"
    with: "db.prod.internal"
  - key: API_URL
    pattern: "http"
    with: "https"
```

```go
rules, err := envreplace.LoadRules("replace-rules.yaml")
```

## Notes

- Input maps are never mutated; a new map is always returned.
- Rules are applied in order; later rules see the output of earlier ones.
- An empty `pattern` or an invalid regex returns an error immediately.
