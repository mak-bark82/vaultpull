package audit

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
)

func TestNewLogger_NoPath(t *testing.T) {
	l, err := NewLogger("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.file != nil {
		t.Fatal("expected nil file for empty path")
	}
}

func TestLog_WritesEntry(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "audit-*.log")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	l, err := NewLogger(tmp.Name())
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	defer l.Close()

	e := Entry{
		VaultPath: "secret/myapp",
		EnvFile:   ".env",
		Keys:      []string{"DB_PASS", "API_KEY"},
		Status:    "success",
	}
	if err := l.Log(e); err != nil {
		t.Fatalf("Log: %v", err)
	}
	l.Close()

	f, _ := os.Open(tmp.Name())
	defer f.Close()
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("expected one log line")
	}
	var got Entry
	if err := json.Unmarshal(scanner.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.VaultPath != e.VaultPath {
		t.Errorf("vault_path: got %q want %q", got.VaultPath, e.VaultPath)
	}
	if got.Status != "success" {
		t.Errorf("status: got %q want success", got.Status)
	}
	if got.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLog_NoOp_NilFile(t *testing.T) {
	l := &Logger{}
	if err := l.Log(Entry{Status: "ok"}); err != nil {
		t.Fatalf("unexpected error on no-op log: %v", err)
	}
}
