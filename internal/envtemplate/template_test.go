package envtemplate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/envtemplate"
)

func writeTemplate(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env.template")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write template: %v", err)
	}
	return p
}

func TestParseTemplate_Valid(t *testing.T) {
	p := writeTemplate(t, "# comment\nDB_HOST\nDB_PORT=5432\nSECRET_KEY\n")
	entries, err := envtemplate.ParseTemplate(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[1].Key != "DB_PORT" || !entries[1].HasDefault || entries[1].Default != "5432" {
		t.Errorf("unexpected entry: %+v", entries[1])
	}
	if entries[0].HasDefault {
		t.Errorf("DB_HOST should have no default")
	}
}

func TestParseTemplate_NotFound(t *testing.T) {
	_, err := envtemplate.ParseTemplate("/nonexistent/.env.template")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseTemplate_SkipsEmpty(t *testing.T) {
	p := writeTemplate(t, "\n# only comments\n\n")
	entries, err := envtemplate.ParseTemplate(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestApplyDefaults_FillsMissing(t *testing.T) {
	entries := []envtemplate.Entry{
		{Key: "DB_HOST", Default: "localhost", HasDefault: true},
		{Key: "DB_PORT", Default: "5432", HasDefault: true},
		{Key: "SECRET", HasDefault: false},
	}
	secrets := map[string]string{"DB_HOST": "prod-db"}
	result := envtemplate.ApplyDefaults(secrets, entries)
	if result["DB_HOST"] != "prod-db" {
		t.Errorf("should not overwrite existing key")
	}
	if result["DB_PORT"] != "5432" {
		t.Errorf("should fill default for DB_PORT")
	}
	if _, ok := result["SECRET"]; ok {
		t.Errorf("SECRET has no default and was not in secrets; should be absent")
	}
}
