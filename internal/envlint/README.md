# envlint

The `envlint` package provides a configurable linter for environment variable secrets.

It checks keys and values against a set of rules and reports violations without modifying the input.

## Default Rules

| Rule | Description |
|---|---|
| `no-empty-value` | Value must not be blank or whitespace-only |
| `key-uppercase` | Key must be fully uppercase |
| `no-spaces-in-key` | Key must not contain spaces |
| `valid-key-chars` | Key must contain only `A-Z`, `0-9`, and `_` |

## Usage

```go
linter := envlint.New()
violations := linter.Lint(secrets)
for _, v := range violations {
    fmt.Println(v.Error())
}
```

## Custom Rules

```go
rules := []envlint.Rule{
    {
        Name:    "no-localhost",
        Message: "value must not be localhost in production",
        Check:   func(_, value string) bool { return value == "localhost" },
    },
}
linter := envlint.WithRules(rules)
```

## Violation

Each `Violation` exposes `Key`, `Rule`, `Message`, and implements `error` via `.Error()`.
