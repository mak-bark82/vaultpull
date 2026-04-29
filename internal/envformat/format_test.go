package envformat

import (
	"strings"
	"testing"
)

var baseEnv = map[string]string{
	"DB_HOST": "localhost",
	"DB_PORT": "5432",
	"APP_ENV": "production",
}

func TestFormat_PlainStyle(t *testing.T) {
	opts := DefaultOptions()
	out := Format(baseEnv, opts)

	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected plain key=value, got: %s", out)
	}
	if strings.Contains(out, "export ") {
		t.Errorf("plain style should not include 'export', got: %s", out)
	}
}

func TestFormat_ExportStyle(t *testing.T) {
	opts := DefaultOptions()
	opts.Style = StyleExport
	out := Format(baseEnv, opts)

	if !strings.Contains(out, "export DB_HOST=localhost") {
		t.Errorf("expected export prefix, got: %s", out)
	}
	if !strings.Contains(out, "export APP_ENV=production") {
		t.Errorf("expected export prefix for APP_ENV, got: %s", out)
	}
}

func TestFormat_QuotedStyle(t *testing.T) {
	opts := DefaultOptions()
	opts.Style = StyleQuoted
	out := Format(map[string]string{"KEY": "hello world"}, opts)

	if !strings.Contains(out, `KEY="hello world"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestFormat_InlineStyle(t *testing.T) {
	opts := DefaultOptions()
	opts.Style = StyleInline
	opts.Separator = " ; "
	out := Format(map[string]string{"A": "1", "B": "2"}, opts)

	if !strings.Contains(out, " ; ") {
		t.Errorf("expected inline separator, got: %s", out)
	}
	if strings.Contains(out, "\n") {
		t.Errorf("inline style should not contain newlines, got: %s", out)
	}
}

func TestFormat_SortedKeys(t *testing.T) {
	opts := DefaultOptions()
	opts.SortKeys = true
	out := Format(baseEnv, opts)

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "APP_ENV") {
		t.Errorf("expected APP_ENV first (sorted), got: %s", lines[0])
	}
}

func TestFormat_EmptyMap(t *testing.T) {
	out := Format(map[string]string{}, DefaultOptions())
	if out != "" {
		t.Errorf("expected empty output for empty map, got: %q", out)
	}
}

func TestFormat_DefaultSeparatorFallback(t *testing.T) {
	opts := Options{
		Style:     StyleInline,
		SortKeys:  false,
		Separator: "",
	}
	out := Format(map[string]string{"X": "1", "Y": "2"}, opts)
	if !strings.Contains(out, " ; ") {
		t.Errorf("expected default separator fallback, got: %s", out)
	}
}
