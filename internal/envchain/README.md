# envchain

Provides a composable pipeline for applying sequential transformations to a secret map.

## Usage

```go
c := envchain.New()
c.Add("upper", func(m map[string]string) (map[string]string, error) {
    out := make(map[string]string)
    for k, v := range m {
        out[k] = strings.ToUpper(v)
    }
    return out, nil
})

results, final, err := c.Run(secrets)
```

## Concepts

- **Stage** — a named function that transforms a `map[string]string`.
- **Chain** — an ordered list of stages.
- **Result** — per-stage snapshot of before/after values and any error.

## Behaviour

- Stages are applied in the order they are added.
- If a stage returns an error, the chain halts and the error is returned.
- The original input map is never mutated.
- Each stage receives a copy of the previous stage's output.
