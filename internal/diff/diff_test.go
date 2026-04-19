package diff

import (
	"testing"
)

func TestCompare_Added(t *testing.T) {
	existing := map[string]string{}
	incoming := map[string]string{"FOO": "bar"}
	r := Compare(existing, incoming)
	if len(r.Added) != 1 || r.Added["FOO"] != "bar" {
		t.Errorf("expected FOO to be added")
	}
	if !r.HasChanges() {
		t.Error("expected HasChanges to be true")
	}
}

func TestCompare_Changed(t *testing.T) {
	existing := map[string]string{"FOO": "old"}
	incoming := map[string]string{"FOO": "new"}
	r := Compare(existing, incoming)
	if len(r.Changed) != 1 || r.Changed["FOO"] != "new" {
		t.Errorf("expected FOO to be changed")
	}
}

func TestCompare_Removed(t *testing.T) {
	existing := map[string]string{"FOO": "bar", "BAZ": "qux"}
	incoming := map[string]string{"FOO": "bar"}
	r := Compare(existing, incoming)
	if len(r.Removed) != 1 {
		t.Errorf("expected BAZ to be removed")
	}
}

func TestCompare_Unchanged(t *testing.T) {
	existing := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"FOO": "bar"}
	r := Compare(existing, incoming)
	if len(r.Unchanged) != 1 {
		t.Errorf("expected FOO to be unchanged")
	}
	if r.HasChanges() {
		t.Error("expected HasChanges to be false")
	}
}

func TestSummary(t *testing.T) {
	r := Result{
		Added:     map[string]string{"A": "1"},
		Changed:   map[string]string{},
		Removed:   map[string]string{},
		Unchanged: map[string]string{"B": "2"},
	}
	s := r.Summary()
	if s != "added=1 changed=0 removed=0 unchanged=1" {
		t.Errorf("unexpected summary: %s", s)
	}
}
