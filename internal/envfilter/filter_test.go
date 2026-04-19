package envfilter

import (
	"testing"
)

func baseSecrets() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"API_KEY":     "abc123",
		"DEBUG":       "true",
	}
}

func TestFilter_NoRules(t *testing.T) {
	f := &Filter{}
	out := f.Apply(baseSecrets())
	if len(out) != 4 {
		t.Errorf("expected 4 keys, got %d", len(out))
	}
}

func TestFilter_IncludeOnly(t *testing.T) {
	f := &Filter{Include: []string{"DB_HOST", "API_KEY"}}
	out := f.Apply(baseSecrets())
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been excluded")
	}
}

func TestFilter_ExcludeOnly(t *testing.T) {
	f := &Filter{Exclude: []string{"DEBUG", "DB_PASSWORD"}}
	out := f.Apply(baseSecrets())
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DEBUG"]; ok {
		t.Error("DEBUG should have been excluded")
	}
}

func TestFilter_IncludeAndExclude(t *testing.T) {
	f := &Filter{
		Include: []string{"DB_HOST", "DB_PASSWORD", "API_KEY"},
		Exclude: []string{"DB_PASSWORD"},
	}
	out := f.Apply(baseSecrets())
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been excluded")
	}
}

func TestFilter_EmptySecrets(t *testing.T) {
	f := &Filter{Include: []string{"DB_HOST"}}
	out := f.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected 0 keys, got %d", len(out))
	}
}
