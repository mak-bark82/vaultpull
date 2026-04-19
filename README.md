# vaultpull

> CLI tool to sync secrets from HashiCorp Vault into local `.env` files safely

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

Authenticate with Vault and pull secrets into a `.env` file:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your-vault-token"

vaultpull --path secret/data/myapp --output .env
```

This will fetch all key-value pairs from the specified Vault path and write them to `.env` in the format:

```
KEY=value
ANOTHER_KEY=another_value
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path to read from | *(required)* |
| `--output` | Output file path | `.env` |
| `--overwrite` | Overwrite existing file | `false` |
| `--addr` | Vault server address | `$VAULT_ADDR` |

### Example with overwrite

```bash
vaultpull --path secret/data/prod --output .env.production --overwrite
```

---

## Requirements

- Go 1.21+
- A running HashiCorp Vault instance
- A valid Vault token or supported auth method

---

## License

[MIT](LICENSE) © 2024 yourusername