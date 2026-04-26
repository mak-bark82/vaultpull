package envpriority_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/envpriority"
)

func baseSources() []envpriority.Source {
	return []envpriority.Source{
		{Name: "vault", Priority: 1, Values: map[string]string{"DB_HOST": "vault-host", "DB_PORT": "5432"}},
		{Name: "local", Priority: 2, Values: map[string]string{"DB_HOST": "local-host", "APP_ENV": "dev"}},
	}
}

func TestMerge_HigherPriorityWins(t *testing.T) {
	r, err := envpriority.Merge(baseSources())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := r.Merged["DB_HOST"]; got != "vault-host" {
		t.Errorf("expected vault-host, got %s", got)
	}
	if r.Origin["DB_HOST"] != "vault" {
		t.Errorf("expected origin vault, got %s", r.Origin["DB_HOST"])
	}
}

func TestMerge_LowerPriorityFillsMissingKeys(t *testing.T) {
	r, err := envpriority.Merge(baseSources())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := r.Merged["APP_ENV"]; got != "dev" {
		t.Errorf("expected dev, got %s", got)
	}
}

func TestMerge_EmptySourceName_ReturnsError(t *testing.T) {
	srcs := []envpriority.Source{{Name: "", Priority: 1, Values: map[string]string{"K": "v"}}}
	_, err := envpriority.Merge(srcs)
	if err == nil {
		t.Fatal("expected error for empty source name")
	}
}

func TestMerge_InvalidPriority_ReturnsError(t *testing.T) {
	srcs := []envpriority.Source{{Name: "x", Priority: 0, Values: map[string]string{}}}
	_, err := envpriority.Merge(srcs)
	if err == nil {
		t.Fatal("expected error for priority 0")
	}
}

func TestMerge_EqualPriority_FirstSourceWins(t *testing.T) {
	srcs := []envpriority.Source{
		{Name: "a", Priority: 2, Values: map[string]string{"KEY": "from-a"}},
		{Name: "b", Priority: 2, Values: map[string]string{"KEY": "from-b"}},
	}
	r, err := envpriority.Merge(srcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Merged["KEY"] != "from-a" {
		t.Errorf("expected from-a, got %s", r.Merged["KEY"])
	}
}

func TestSummary_CountsPerSource(t *testing.T) {
	r, _ := envpriority.Merge(baseSources())
	counts := envpriority.Summary(r)
	if counts["vault"] != 2 {
		t.Errorf("expected vault=2, got %d", counts["vault"])
	}
	if counts["local"] != 1 {
		t.Errorf("expected local=1, got %d", counts["local"])
	}
}

func TestMerge_EmptySources_ReturnsEmptyResult(t *testing.T) {
	r, err := envpriority.Merge(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Merged) != 0 {
		t.Errorf("expected empty merged map")
	}
}
