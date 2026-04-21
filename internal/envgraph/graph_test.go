package envgraph

import (
	"testing"
)

func TestResolve_NoDeps(t *testing.T) {
	g := New()
	g.Add("A", nil)
	g.Add("B", nil)

	order, err := g.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(order))
	}
}

func TestResolve_DepsFirst(t *testing.T) {
	g := New()
	g.Add("A", []string{"B"})
	g.Add("B", nil)

	order, err := g.Resolve()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	idxA, idxB := -1, -1
	for i, k := range order {
		if k == "A" {
			idxA = i
		}
		if k == "B" {
			idxB = i
		}
	}
	if idxB >= idxA {
		t.Errorf("expected B before A, got order: %v", order)
	}
}

func TestResolve_CycleDetected(t *testing.T) {
	g := New()
	g.Add("A", []string{"B"})
	g.Add("B", []string{"A"})

	_, err := g.Resolve()
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestBuildFromEnv_ExtractsRefs(t *testing.T) {
	env := map[string]string{
		"DB_URL": "postgres://${DB_USER}:${DB_PASS}@localhost",
		"DB_USER": "admin",
		"DB_PASS": "secret",
	}
	g := BuildFromEnv(env)

	node, ok := g.nodes["DB_URL"]
	if !ok {
		t.Fatal("expected DB_URL node")
	}
	if len(node.Deps) != 2 {
		t.Fatalf("expected 2 deps for DB_URL, got %d", len(node.Deps))
	}
}

func TestBuildFromEnv_NoRefs(t *testing.T) {
	env := map[string]string{
		"PLAIN": "value",
	}
	g := BuildFromEnv(env)
	node := g.nodes["PLAIN"]
	if len(node.Deps) != 0 {
		t.Errorf("expected no deps, got %v", node.Deps)
	}
}

func TestExtractRefs_Multiple(t *testing.T) {
	refs := extractRefs("${A} and ${B} and ${C}")
	if len(refs) != 3 {
		t.Fatalf("expected 3 refs, got %d: %v", len(refs), refs)
	}
}

func TestExtractRefs_None(t *testing.T) {
	refs := extractRefs("no references here")
	if len(refs) != 0 {
		t.Errorf("expected no refs, got %v", refs)
	}
}
