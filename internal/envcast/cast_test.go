package envcast_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/envcast"
)

func TestCast_StringPassthrough(t *testing.T) {
	secrets := map[string]string{"HOST": "localhost"}
	results := envcast.Cast(secrets, envcast.Options{Types: map[string]string{}})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Typed != "localhost" {
		t.Errorf("expected 'localhost', got %v", results[0].Typed)
	}
}

func TestCast_Bool(t *testing.T) {
	secrets := map[string]string{"ENABLED": "true"}
	results := envcast.Cast(secrets, envcast.Options{Types: map[string]string{"ENABLED": "bool"}})
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if results[0].Typed != true {
		t.Errorf("expected true, got %v", results[0].Typed)
	}
}

func TestCast_Int(t *testing.T) {
	secrets := map[string]string{"PORT": "8080"}
	results := envcast.Cast(secrets, envcast.Options{Types: map[string]string{"PORT": "int"}})
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if results[0].Typed != int64(8080) {
		t.Errorf("expected 8080, got %v", results[0].Typed)
	}
}

func TestCast_Float(t *testing.T) {
	secrets := map[string]string{"RATIO": "3.14"}
	results := envcast.Cast(secrets, envcast.Options{Types: map[string]string{"RATIO": "float"}})
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if results[0].Typed.(float64) < 3.13 {
		t.Errorf("expected ~3.14, got %v", results[0].Typed)
	}
}

func TestCast_InvalidBool(t *testing.T) {
	secrets := map[string]string{"ENABLED": "notabool"}
	results := envcast.Cast(secrets, envcast.Options{Types: map[string]string{"ENABLED": "bool"}})
	if results[0].Err == nil {
		t.Error("expected error for invalid bool")
	}
}

func TestCast_InvalidInt(t *testing.T) {
	secrets := map[string]string{"PORT": "abc"}
	results := envcast.Cast(secrets, envcast.Options{Types: map[string]string{"PORT": "int"}})
	if results[0].Err == nil {
		t.Error("expected error for invalid int")
	}
}

func TestCast_UnknownType(t *testing.T) {
	secrets := map[string]string{"X": "val"}
	results := envcast.Cast(secrets, envcast.Options{Types: map[string]string{"X": "uuid"}})
	if results[0].Err == nil {
		t.Error("expected error for unknown type")
	}
}

func TestCastOne_String(t *testing.T) {
	r := envcast.CastOne("KEY", "hello", "string")
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if r.Typed != "hello" {
		t.Errorf("expected 'hello', got %v", r.Typed)
	}
}

func TestCastOne_InvalidFloat(t *testing.T) {
	r := envcast.CastOne("RATIO", "not-a-float", "float")
	if r.Err == nil {
		t.Error("expected error for invalid float")
	}
}
