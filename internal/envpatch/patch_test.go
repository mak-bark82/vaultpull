package envpatch_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envpatch"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "secret",
	}
}

func TestApply_SetAddsOrUpdates(t *testing.T) {
	patches := []envpatch.Patch{
		{Op: envpatch.OpSet, Key: "NEW_KEY", Value: "new_value"},
		{Op: envpatch.OpSet, Key: "DB_HOST", Value: "remotehost"},
	}
	out, result, err := envpatch.Apply(baseSecrets(), patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "new_value" {
		t.Errorf("expected NEW_KEY=new_value, got %q", out["NEW_KEY"])
	}
	if out["DB_HOST"] != "remotehost" {
		t.Errorf("expected DB_HOST=remotehost, got %q", out["DB_HOST"])
	}
	if len(result.Applied) != 2 {
		t.Errorf("expected 2 applied, got %d", len(result.Applied))
	}
}

func TestApply_DeleteRemovesKey(t *testing.T) {
	patches := []envpatch.Patch{
		{Op: envpatch.OpDelete, Key: "API_KEY"},
	}
	out, result, err := envpatch.Apply(baseSecrets(), patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := out["API_KEY"]; exists {
		t.Error("expected API_KEY to be deleted")
	}
	if len(result.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(result.Applied))
	}
}

func TestApply_Delete_SkipsMissingKey(t *testing.T) {
	patches := []envpatch.Patch{
		{Op: envpatch.OpDelete, Key: "NONEXISTENT"},
	}
	_, result, err := envpatch.Apply(baseSecrets(), patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "NONEXISTENT" {
		t.Errorf("expected NONEXISTENT in skipped, got %v", result.Skipped)
	}
}

func TestApply_RenameMovesKey(t *testing.T) {
	patches := []envpatch.Patch{
		{Op: envpatch.OpRename, Key: "DB_PORT", NewKey: "DATABASE_PORT"},
	}
	out, result, err := envpatch.Apply(baseSecrets(), patches)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := out["DB_PORT"]; exists {
		t.Error("expected DB_PORT to be removed after rename")
	}
	if out["DATABASE_PORT"] != "5432" {
		t.Errorf("expected DATABASE_PORT=5432, got %q", out["DATABASE_PORT"])
	}
	if len(result.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(result.Applied))
	}
}

func TestApply_UnknownOp_ReturnsError(t *testing.T) {
	patches := []envpatch.Patch{
		{Op: "upsert", Key: "FOO"},
	}
	_, _, err := envpatch.Apply(baseSecrets(), patches)
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestApply_MissingKey_ReturnsError(t *testing.T) {
	patches := []envpatch.Patch{
		{Op: envpatch.OpSet, Key: "", Value: "val"},
	}
	_, _, err := envpatch.Apply(baseSecrets(), patches)
	if err == nil {
		t.Error("expected error for missing key field")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	src := baseSecrets()
	original := src["DB_HOST"]
	patches := []envpatch.Patch{
		{Op: envpatch.OpSet, Key: "DB_HOST", Value: "changed"},
	}
	_, _, _ = envpatch.Apply(src, patches)
	if src["DB_HOST"] != original {
		t.Error("Apply mutated the input map")
	}
}
