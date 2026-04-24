# envaccess

Provides rule-based access control for secret keys pulled from Vault.

## Overview

`envaccess` lets you define which keys a given run is permitted to read or write, preventing accidental exposure of sensitive secrets.

## Rule format (`access.yaml`)

```yaml
rules:
  - pattern: DB_*
    permission: write
  - pattern: API_KEY
    permission: read
  - pattern: SECRET_*
    permission: none
```

### Permissions

| Value   | Meaning                        |
|---------|--------------------------------|
| `none`  | Key is blocked entirely        |
| `read`  | Key may be read                |
| `write` | Key may be read and written    |

## Pattern matching

Patterns support a single trailing `*` wildcard (prefix match). Exact patterns are also supported. **The first matching rule wins.**

## Usage

```go
rules, err := envaccess.LoadRules("access.yaml")
checker, err := envaccess.New(rules)

// Check a key
perm := checker.Check("DB_HOST") // envaccess.PermWrite

// Enforce a minimum permission (returns error if insufficient)
if err := checker.Enforce("API_KEY", envaccess.PermWrite); err != nil {
    log.Fatal(err)
}
```
