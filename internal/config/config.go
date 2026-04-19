package config

import (
	"errors"
	"os"

	"github.com/jo)

// Config holds the configuration for vaultpull.
type Config struct {
	VaultAddress string
	VaultToken   string
	SecretPath   string
	OutputFile   string
	Overwrite    bool
}

// Load reads configuration from environment variables,
// optionally loading a .env file first if it exists.
func Load(envFile string) (*Config, error) {
	if envFile != "" {
		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Load(envFile); err != nil {
				return nil, err
			}
		}
	}

	cfg := &Config{
		VaultAddress: os.Getenv("VAULT_ADDR"),
		VaultToken:   os.Getenv("VAULT_TOKEN"),
		SecretPath:   os.Getenv("VAULTPULL_SECRET_PATH"),
		OutputFile:   os.Getenv("VAULTPULL_OUTPUT_FILE"),
		Overwrite:    os.Getenv("VAULTPULL_OVERWRITE") == "true",
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	if cfg.OutputFile == "" {
		cfg.OutputFile = ".env"
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.VaultAddress == "" {
		return errors.New("VAULT_ADDR is required")
	}
	if c.VaultToken == "" {
		return errors.New("VAULT_TOKEN is required")
	}
	if c.SecretPath == "" {
		return errors.New("VAULTPULL_SECRET_PATH is required")
	}
	return nil
}
