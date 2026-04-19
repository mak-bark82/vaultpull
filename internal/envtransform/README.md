# envtransform

The `envtransform` package provides lightweight value transformation for secret maps pulled from Vault before they are written to `.env` files.

## Usage

```go
rule := envtransform.Rule{
    TrimSpace: true,
    UpperCase: true,
    Prefix:    "APP_",
}

transformed := envtransform.Transform(secrets, rule)
```

## Rule fields

| Field       | Type   | Description                              |
|-------------|--------|------------------------------------------|
| `Prefix`    | string | Prepend string to every value            |
| `Suffix`    | string | Append string to every value             |
| `UpperCase` | bool   | Convert value to upper-case              |
| `LowerCase` | bool   | Convert value to lower-case              |
| `TrimSpace` | bool   | Trim leading/trailing whitespace         |

`UpperCase` takes precedence over `LowerCase` if both are set.

Transformations are applied in order: TrimSpace → Case → Prefix → Suffix.
