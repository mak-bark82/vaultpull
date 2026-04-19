package syncer_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/syncer"
)

func mockVaultServer(t *testing.T, body string, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func TestRun_Success(t *testing.T) {
	body := `{"data":{"data":{"KEY":"value","SECRET":"abc"}}}`
	srv := mockVaultServer(t, body, http.StatusOK)
	defer srv.Close()

	dir := t.TempDir()
	outFile := filepath.Join(dir, ".env")

	cfg := &config.Config{
		VaultAddr:  srv.URL,
		VaultToken: "test-token",
		Mappings: []config.Mapping{
			{VaultPath: "secret/data/app", EnvFile: outFile, Overwrite: true},
		},
	}

	s, err := syncer.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results := s.Run()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Fatal("expected env file to be created")
	}
}

func TestRun_VaultError(t *testing.T) {
	srv := mockVaultServer(t, `{"errors":["permission denied"]}`, http.StatusForbidden)
	defer srv.Close()

	dir := t.TempDir()
	cfg := &config.Config{
		VaultAddr:  srv.URL,
		VaultToken: "bad-token",
		Mappings: []config.Mapping{
			{VaultPath: "secret/data/app", EnvFile: filepath.Join(dir, ".env"), Overwrite: true},
		},
	}

	s, err := syncer.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results := s.Run()
	if results[0].Err == nil {
		t.Fatal("expected error, got nil")
	}
}
