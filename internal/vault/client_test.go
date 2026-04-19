package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_MissingAddress(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("VAULT_TOKEN", "test-token")

	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected error for missing address, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", "")

	_, err := NewClient(Config{Address: "http://127.0.0.1:8200"})
	if err == nil {
		t.Fatal("expected error for missing token, got nil")
	}
}

func TestReadSecrets_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"data":{"DB_HOST":"localhost","DB_PORT":"5432"}}}`))
	}))
	defer server.Close()

	client, err := NewClient(Config{Address: server.URL, Token: "fake-token"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	secrets, err := client.ReadSecrets("secret", "myapp")
	if err != nil {
		t.Fatalf("unexpected error reading secrets: %v", err)
	}

	if secrets["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", secrets["DB_HOST"])
	}
	if secrets["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", secrets["DB_PORT"])
	}
}

func TestReadSecrets_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := NewClient(Config{Address: server.URL, Token: "fake-token"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.ReadSecrets("secret", "missing")
	if err == nil {
		t.Fatal("expected error for missing secret, got nil")
	}
}
