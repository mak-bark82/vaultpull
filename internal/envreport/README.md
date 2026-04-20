# envreport

Package `envreport` generates human-readable sync reports after a vaultpull operation.

## Overview

After secrets are pulled from Vault and merged into a local `.env` file, `envreport` collects per-key status information and renders a formatted summary.

## Usage

```go
import "github.com/your-org/vaultpull/internal/envreport"

entries := []envreport.Entry{
    {Key: "DB_HOST", Status: "added",   Source: "secret/myapp"},
    {Key: "DB_PASS", Status: "changed", Source: "secret/myapp"},
    {Key: "APP_ENV", Status: "unchanged", Source: "secret/myapp"},
}

r := envreport.New(".env", entries)
r.Render(os.Stdout)
```

## Entry Statuses

| Status      | Meaning                                      |
|-------------|----------------------------------------------|
| `added`     | Key did not exist locally; written from Vault |
| `changed`   | Key existed but value was updated             |
| `removed`   | Key no longer present in Vault               |
| `unchanged` | Key exists with identical value              |

## Output Example

```
vaultpull sync report
Timestamp : 2024-05-01T12:00:00Z
Env file  : .env
------------------------------------------------
  [added    ] DB_HOST  (secret/myapp)
  [changed  ] DB_PASS  (secret/myapp)
  [unchanged] APP_ENV  (secret/myapp)
------------------------------------------------
added: 1  changed: 1  removed: 0  unchanged: 1
```
