# envcompare

Package `envcompare` provides utilities for comparing two sets of environment
variables — either in-memory maps or `.env` files on disk.

## Features

- **Compare** two `map[string]string` values and categorise every key into one
  of four buckets: `OnlyInLeft`, `OnlyInRight`, `Different`, or `Identical`.
- **Render** a human-readable diff to any `io.Writer` using familiar `<`, `>`,
  and `~` prefix symbols.
- **Summary** returns a compact one-line description of the result (e.g.
  `"1 only-left, 2 changed, 5 identical"`).
- **LoadFile** parses a `.env` file, stripping comments, blank lines, and
  surrounding quotes.
- **CompareFiles** is a convenience wrapper that loads two files and returns a
  `Result` directly.

## Usage

```go
result, err := envcompare.CompareFiles("staging.env", "production.env")
if err != nil {
    log.Fatal(err)
}
envcompare.Render(os.Stdout, result, false)
fmt.Println(envcompare.Summary(result))
```

Pass `redact = true` to `Render` to replace all values with `***` — useful
when printing to logs that may be captured in CI.
