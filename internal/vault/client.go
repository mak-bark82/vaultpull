package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client.
type Client struct {
	vc *vaultapi.Client
}

// Config holds connection parameters for Vault.
type Config struct {
	Address string
	Token   string
}

// NewClient creates and configures a new Vault client.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()

	addr := cfg.Address
	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
	}
	if addr == "" {
		return nil, fmt.Errorf("vault address is required (set --vault-addr or VAULT_ADDR)")
	}
	vcfg.Address = addr

	client, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token is required (set --vault-token or VAULT_TOKEN)")
	}
	client.SetToken(token)

	return &Client{vc: client}, nil
}

// ReadSecrets reads KV v2 secrets at the given mount and path.
func (c *Client) ReadSecrets(mount, path string) (map[string]string, error) {
	fullPath := fmt.Sprintf("%s/data/%s", mount, path)
	secret, err := c.vc.Logical().Read(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %q: %w", fullPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret found at %q", fullPath)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected secret format at %q", fullPath)
	}

	result := make(map[string]string, len(data))
	for k, v := range data {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}
