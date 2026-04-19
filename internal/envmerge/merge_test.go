package envmerge_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/envmerge"
)

func TestMerge_AddsNewKeys(t *testing.T) {
	local := map[string]string{"EXISTING": "old"}
	incoming := map[string]string{"NEW_KEY": "value"}

	res := envmerge.Merge(local, incoming, envmerge.PreferVault)

	if res.Merged["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %q", res.Merged["NEW_KEY"])
	}
	if res.Added != 1 {
		t.Errorf("expected Added=1, got %d", res.Added)
	}
}

func TestMerge_PreferVault_Overwrites(t *testing.T) {
	local := map[string]string{"KEY": "local_val"}
	incoming := map[string]string{"KEY": "vault_val"}

	res := envmerge.Merge(local, incoming, envmerge.PreferVault)

	if res.Merged["KEY"] != "vault_val" {
		t.Errorf("expected vault_val, got %q", res.Merged["KEY"])
	}
	if res.Overwritten != 1 {
		t.Errorf("expected Overwritten=1, got %d", res.Overwritten)
	}
}

func TestMerge_PreferLocal_Skips(t *testing.T) {
	local := map[string]string{"KEY": "local_val"}
	incoming := map[string]string{"KEY": "vault_val"}

	res := envmerge.Merge(local, incoming, envmerge.PreferLocal)

	if res.Merged["KEY"] != "local_val" {
		t.Errorf("expected local_val to be preserved, got %q", res.Merged["KEY"])
	}
	if res.Skipped != 1 {
		t.Errorf("expected Skipped=1, got %d", res.Skipped)
	}
}

func TestMerge_UnchangedValue_NotCounted(t *testing.T) {
	local := map[string]string{"KEY": "same"}
	incoming := map[string]string{"KEY": "same"}

	res := envmerge.Merge(local, incoming, envmerge.PreferVault)

	if res.Overwritten != 0 || res.Skipped != 0 || res.Added != 0 {
		t.Errorf("expected no changes, got overwritten=%d skipped=%d added=%d",
			res.Overwritten, res.Skipped, res.Added)
	}
}

func TestMerge_EmptyLocal(t *testing.T) {
	local := map[string]string{}
	incoming := map[string]string{"A": "1", "B": "2"}

	res := envmerge.Merge(local, incoming, envmerge.PreferVault)

	if res.Added != 2 {
		t.Errorf("expected Added=2, got %d", res.Added)
	}
	if len(res.Merged) != 2 {
		t.Errorf("expected 2 merged keys, got %d", len(res.Merged))
	}
}

func TestMerge_LocalKeysPreserved(t *testing.T) {
	// Keys present only in local should always be retained in the merged output.
	local := map[string]string{"LOCAL_ONLY": "keep", "SHARED": "local_val"}
	incoming := map[string]string{"SHARED": "vault_val"}

	res := envmerge.Merge(local, incoming, envmerge.PreferVault)

	if res.Merged["LOCAL_ONLY"] != "keep" {
		t.Errorf("expected LOCAL_ONLY to be preserved, got %q", res.Merged["LOCAL_ONLY"])
	}
}
