# envnamespace

Provides namespace-aware key qualification and stripping for environment variable maps.

## Overview

A **namespace** pairs a logical name (e.g. `prod`) with a key prefix (e.g. `PROD_`).  
The `Resolver` uses these definitions to:

- **Qualify** — add the namespace prefix to every key in a map.
- **Strip** — remove the namespace prefix from keys that carry it; pass others through unchanged.

## Usage

```go
ns := []envnamespace.Namespace{
    {Name: "prod",    Prefix: "PROD_"},
    {Name: "staging", Prefix: "STG_"},
}

r, err := envnamespace.NewResolver(ns)
if err != nil {
    log.Fatal(err)
}

// Add prefix before writing to a shared store.
qualified, err := r.Qualify("prod", secrets)

// Remove prefix after reading from a shared store.
bare, err := r.Strip("prod", qualified)
```

## Error handling

- `NewResolver` returns an error if any namespace has an empty `Name` or `Prefix`.
- `Qualify` and `Strip` return an error when the requested namespace name is not registered.
