package main

import (
	"fmt"
	"os"

	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/syncer"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	mappingsFile := os.Getenv("VAULTPULL_MAPPINGS")
	mappings, err := config.LoadMappings(mappingsFile)
	if err != nil {
		return fmt.Errorf("loading mappings: %w", err)
	}
	if len(mappings) == 0 {
		return fmt.Errorf("no mappings defined; set VAULTPULL_MAPPINGS to a JSON file")
	}
	cfg.Mappings = mappings

	s, err := syncer.New(cfg)
	if err != nil {
		return err
	}

	results := s.Run()
	exitCode := 0
	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "FAIL  %s -> %s: %v\n", r.Path, r.OutFile, r.Err)
			exitCode = 1
		} else if r.Skipped {
			fmt.Printf("SKIP  %s -> %s (already exists)\n", r.Path, r.OutFile)
		} else {
			fmt.Printf("OK    %s -> %s\n", r.Path, r.OutFile)
		}
	}

	if exitCode != 0 {
		os.Exit(exitCode)
	}
	return nil
}
