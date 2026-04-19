package envexpand

import (
	"os"
	"testing"
)

func TestExpand_ResolvesInternalReference(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api",
	}
	result := Expand(env)
	if result["API_URL"] != "https://example.com/api" {
		t.Errorf("expected expanded URL, got %q", result["API_URL"])
	}
}

func TestExpand_FallsBackToOS(t *testing.T) {
	os.Setenv("OS_VAR", "from-os")
	defer os.Unsetenv("OS_VAR")

	env := map[string]string{
		"VALUE": "${OS_VAR}-suffix",
	}
	result := Expand(env)
	if result["VALUE"] != "from-os-suffix" {
		t.Errorf("expected os fallback, got %q", result["VALUE"])
	}
}

func TestExpand_NoReferences(t *testing.T) {
	env := map[string]string{
		"PLAIN": "hello",
	}
	result := Expand(env)
	if result["PLAIN"] != "hello" {
		t.Errorf("expected unchanged value, got %q", result["PLAIN"])
	}
}

func TestExpand_UnresolvedReference(t *testing.T) {
	env := map[string]string{
		"VALUE": "${MISSING_VAR}",
	}
	result := Expand(env)
	if result["VALUE"] != "" {
		t.Errorf("expected empty string for unresolved ref, got %q", result["VALUE"])
	}
}

func TestExpandValue_Direct(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	out := ExpandValue("$HOST:8080", env)
	if out != "localhost:8080" {
		t.Errorf("expected localhost:8080, got %q", out)
	}
}

func TestHasReferences(t *testing.T) {
	if !HasReferences("${FOO}") {
		t.Error("expected true for value with reference")
	}
	if HasReferences("plainvalue") {
		t.Error("expected false for plain value")
	}
}
