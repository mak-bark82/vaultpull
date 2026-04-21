# envimport

Provides utilities for importing key-value pairs from external files into an
existing secrets map.

## Supported formats

| Format     | Constant            | Notes                              |
|------------|---------------------|------------------------------------|
| `.env`     | `FormatDotenv`      | `KEY=VALUE`, ignores `#` comments  |

## Usage

```go
import "github.com/yourusername/vaultpull/internal/envimport"

existing := map[string]string{"KEEP": "me"}

result, err := envimport.Import(".env.staging", existing, envimport.Options{
    Format:    envimport.FormatDotenv,
    Overwrite: false, // preserve keys already present in existing
})
```

## Options

| Field       | Type     | Description                                              |
|-------------|----------|----------------------------------------------------------|
| `Format`    | `Format` | File format to parse (default: `FormatDotenv`)           |
| `Overwrite` | `bool`   | When `true`, incoming values replace existing ones       |

## Notes

- Quoted values (`KEY="hello world"`) are unquoted automatically.
- Malformed lines (no `=` separator) are silently skipped.
- The original `dst` map is never mutated; a new map is returned.
