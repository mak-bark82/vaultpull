# envpolicy

The `envpolicy` package enforces configurable rules against environment variable keys and values before they are written to disk.

## Overview

Rules are defined in a YAML file and evaluated by an `Enforcer`. Each rule specifies:

- `name` — human-readable rule identifier
- `pattern` — regular expression to match against the target
- `target` — `key` (default) or `value`
- `action` — `allow`, `deny`, or `warn`

## Example policy file

```yaml
rules:
  - name: no-debug-keys
    pattern: "^DEBUG"
    target: key
    action: deny

  - name: warn-plaintext-password
    pattern: "password"
    target: value
    action: warn
```

## Usage

```go
rules, err := envpolicy.LoadRules("policy.yaml")
if err != nil {
    log.Fatal(err)
}

enforcer, err := envpolicy.New(rules)
if err != nil {
    log.Fatal(err)
}

violations := enforcer.Check(secrets)
if envpolicy.HasDenials(violations) {
    log.Fatal("policy denied one or more secrets")
}
```

## Actions

| Action | Behaviour |
|--------|-----------|
| `allow` | Explicit allow; no violation recorded |
| `warn`  | Violation recorded but does not block |
| `deny`  | Violation recorded; `HasDenials` returns `true` |
