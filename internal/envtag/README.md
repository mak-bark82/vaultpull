# envtag

The `envtag` package provides metadata tagging for secret keys pulled from Vault.

## Overview

When syncing secrets, you may want to annotate certain keys with metadata — for example, marking secrets as belonging to a specific environment (`env:prod`) or tier (`tier:db`).

## Usage

```go
secrets := map[string]string{"DB_PASS": "s3cr3t", "API_KEY": "abc"}
annotations := map[string]string{"DB_PASS": "env:prod,tier:db"}

tagged := envtag.Parse(secrets, annotations)
prodSecrets := envtag.Filter(tagged, "env")
```

## Tag Format

Annotations are strings with comma-separated `key:value` pairs:

```
env:prod,tier:db,sensitive
```

Tags without a `:` are stored with an empty value.
