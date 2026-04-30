package envnormalize_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/envnormalize"
)

func TestNormalize_UppercaseKeys(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost"}
	res, err := envnormalize.Normalize(secrets, envnormalize.Options{UppercaseKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Secrets["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be present")
	}
	if res.Renamed != 1 {
		t.Errorf("expected 1 renamed, got %d", res.Renamed)
	}
}

func TestNormalize_ReplaceHyphen(t *testing.T) {
	secrets := map[string]string{"my-key": "val"}
	opts := envnormalize.Options{ReplaceHyphen: true}
	res, err := envnormalize.Normalize(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Secrets["my_key"]; !ok {
		t.Error("expected my_key after hyphen replacement")
	}
}

func TestNormalize_StripInvalidChars(t *testing.T) {
	secrets := map[string]string{"key.name!": "v"}
	opts := envnormalize.Options{StripInvalidChars: true}
	res, err := envnormalize.Normalize(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Secrets["keyname"]; !ok {
		t.Errorf("expected 'keyname', got keys: %v", res.Secrets)
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	secrets := map[string]string{"KEY": "  value  "}
	opts := envnormalize.Options{TrimValues: true}
	res, err := envnormalize.Normalize(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Secrets["KEY"] != "value" {
		t.Errorf("expected trimmed value, got %q", res.Secrets["KEY"])
	}
	if res.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", res.Modified)
	}
}

func TestNormalize_KeyCollisionReturnsError(t *testing.T) {
	secrets := map[string]string{"db_host": "a", "DB_HOST": "b"}
	opts := envnormalize.Options{UppercaseKeys: true}
	_, err := envnormalize.Normalize(secrets, opts)
	if err == nil {
		t.Error("expected collision error, got nil")
	}
}

func TestNormalize_EmptyKeyDropped(t *testing.T) {
	secrets := map[string]string{"!!!": "val"}
	opts := envnormalize.Options{StripInvalidChars: true}
	res, err := envnormalize.Normalize(secrets, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Secrets) != 0 {
		t.Errorf("expected empty map after stripping, got %v", res.Secrets)
	}
}

func TestNormalize_DefaultOptions(t *testing.T) {
	secrets := map[string]string{"my-db.host": "  127.0.0.1  "}
	res, err := envnormalize.Normalize(secrets, envnormalize.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := res.Secrets["MY_DBHOST"]; !ok || v != "127.0.0.1" {
		t.Errorf("unexpected result: %v", res.Secrets)
	}
}
