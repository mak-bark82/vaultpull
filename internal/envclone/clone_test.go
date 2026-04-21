package envclone_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/envclone"
)

func baseSrc() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "secret",
	}
}

func TestClone_CopiesNewKeys(t *testing.T) {
	src := baseSrc()
	dst := map[string]string{}

	result, err := envclone.Clone(src, dst, envclone.CloneOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Copied) != 3 {
		t.Errorf("expected 3 copied, got %d", len(result.Copied))
	}
	if dst["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST to be copied")
	}
}

func TestClone_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := baseSrc()
	dst := map[string]string{"DB_HOST": "remotehost"}

	result, err := envclone.Clone(src, dst, envclone.CloneOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if dst["DB_HOST"] != "remotehost" {
		t.Errorf("existing key should not be overwritten")
	}
}

func TestClone_OverwritesWhenEnabled(t *testing.T) {
	src := baseSrc()
	dst := map[string]string{"DB_HOST": "old"}

	result, err := envclone.Clone(src, dst, envclone.CloneOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Overwrite) != 1 {
		t.Errorf("expected 1 overwritten, got %d", len(result.Overwrite))
	}
	if dst["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST to be overwritten")
	}
}

func TestClone_DryRun_DoesNotMutateDst(t *testing.T) {
	src := baseSrc()
	dst := map[string]string{}

	result, err := envclone.Clone(src, dst, envclone.CloneOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Copied) != 3 {
		t.Errorf("expected 3 reported as copied, got %d", len(result.Copied))
	}
	if len(dst) != 0 {
		t.Errorf("dst should remain empty in dry-run mode")
	}
}

func TestClone_NilSrc_ReturnsError(t *testing.T) {
	_, err := envclone.Clone(nil, map[string]string{}, envclone.CloneOptions{})
	if err == nil {
		t.Error("expected error for nil src")
	}
}

func TestClone_NilDst_ReturnsError(t *testing.T) {
	_, err := envclone.Clone(map[string]string{}, nil, envclone.CloneOptions{})
	if err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestSummary_Format(t *testing.T) {
	r := &envclone.Result{
		Copied:    []string{"A", "B"},
		Skipped:   []string{"C"},
		Overwrite: []string{},
	}
	got := r.Summary()
	want := "copied=2 skipped=1 overwritten=0"
	if got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}
