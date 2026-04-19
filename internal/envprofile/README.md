# envprofile

The `envprofile` package provides support for named environment profiles in `vaultpull`.

Profiles allow you to define multiple sync targets (e.g., `dev`, `staging`, `prod`) in a single YAML file, each mapping a Vault prefix to a local `.env` file.

## Profile File Format

```yaml
profiles:
  dev:
    name: dev
    vault_prefix: secret/dev
    env_file: .env.dev
    overrides:
      LOG_LEVEL: debug
  prod:
    name: prod
    vault_prefix: secret/prod
    env_file: .env.prod
```

## Fields

| Field          | Required | Description                              |
|----------------|----------|------------------------------------------|
| `name`         | no       | Human-readable label                     |
| `vault_prefix` | yes      | Vault KV path prefix for this profile    |
| `env_file`     | yes      | Local `.env` file to write secrets into  |
| `overrides`    | no       | Key/value pairs to inject after sync     |

## Usage

```go
ps, err := envprofile.LoadProfiles("profiles.yaml")
p, err := ps.Get("dev")
fmt.Println(p.VaultPrefix) // secret/dev
```
