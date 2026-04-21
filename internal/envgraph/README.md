# envgraph

Package `envgraph` builds and analyzes dependency graphs for environment variables.

## Overview

Env variables often reference each other via `${VAR}` syntax. `envgraph` detects
these relationships, resolves them in topological order, and can export the graph
for visualization.

## Usage

### Build from an env map

```go
env := map[string]string{
    "DB_URL":  "postgres://${DB_USER}:${DB_PASS}@localhost",
    "DB_USER": "admin",
    "DB_PASS": "secret",
}
g := envgraph.BuildFromEnv(env)
```

### Topological resolution

```go
order, err := g.Resolve()
if err != nil {
    log.Fatal(err) // cycle detected
}
fmt.Println(order) // [DB_USER DB_PASS DB_URL]
```

### Export to Graphviz DOT

```go
f, _ := os.Create("graph.dot")
defer f.Close()
g.ExportDOT(f)
// render with: dot -Tpng graph.dot -o graph.png
```

## Cycle Detection

`Resolve` returns an error if a circular dependency is found, e.g.:

```
cycle detected at key: A
```
