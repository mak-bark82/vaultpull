package envgraph

import (
	"strings"
	"testing"
)

func TestExportDOT_ContainsEdges(t *testing.T) {
	g := New()
	g.Add("APP_URL", []string{"BASE_URL"})
	g.Add("BASE_URL", nil)

	var sb strings.Builder
	if err := g.ExportDOT(&sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()

	if !strings.Contains(out, "digraph envgraph") {
		t.Error("expected digraph header")
	}
	if !strings.Contains(out, "BASE_URL") {
		t.Error("expected BASE_URL in output")
	}
	if !strings.Contains(out, "APP_URL") {
		t.Error("expected APP_URL in output")
	}
	if !strings.Contains(out, "->") {
		t.Error("expected edge arrow in output")
	}
}

func TestExportDOT_NoDeps(t *testing.T) {
	g := New()
	g.Add("SOLO", nil)

	var sb strings.Builder
	if err := g.ExportDOT(&sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()

	if !strings.Contains(out, "SOLO") {
		t.Error("expected SOLO node in output")
	}
	if strings.Contains(out, "->") {
		t.Error("expected no edges for isolated node")
	}
}

func TestExportDOT_EmptyGraph(t *testing.T) {
	g := New()

	var sb strings.Builder
	if err := g.ExportDOT(&sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()

	if !strings.Contains(out, "digraph envgraph") {
		t.Error("expected digraph header even for empty graph")
	}
	if !strings.Contains(out, "}") {
		t.Error("expected closing brace")
	}
}

func TestExportDOT_MultipleEdges(t *testing.T) {
	g := New()
	g.Add("C", []string{"A", "B"})
	g.Add("A", nil)
	g.Add("B", nil)

	var sb strings.Builder
	if err := g.ExportDOT(&sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()

	if strings.Count(out, "->") != 2 {
		t.Errorf("expected 2 edges, got:\n%s", out)
	}
}
