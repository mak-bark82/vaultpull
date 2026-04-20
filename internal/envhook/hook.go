// Package envhook provides a mechanism to register and execute lifecycle hooks
// that run before or after a sync operation. Hooks can be used to trigger
// notifications, run validation scripts, or perform custom post-processing.
package envhook

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Phase indicates when a hook should be executed relative to a sync operation.
type Phase string

const (
	// PreSync hooks run before secrets are written to disk.
	PreSync Phase = "pre"
	// PostSync hooks run after secrets have been written to disk.
	PostSync Phase = "post"
)

// Hook represents a single executable hook with a defined phase.
type Hook struct {
	// Phase is either "pre" or "post".
	Phase Phase
	// Command is the shell command to execute.
	Command string
}

// Runner executes registered hooks in order.
type Runner struct {
	hooks  []Hook
	stdout io.Writer
	stderr io.Writer
}

// New creates a new Runner that writes command output to the provided writers.
// If stdout or stderr are nil, os.Stdout and os.Stderr are used respectively.
func New(hooks []Hook, stdout, stderr io.Writer) *Runner {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}
	return &Runner{
		hooks:  hooks,
		stdout: stdout,
		stderr: stderr,
	}
}

// Run executes all hooks matching the given phase in registration order.
// Each hook command is executed via the system shell. If any hook exits with
// a non-zero status, Run returns an error and stops further execution.
func (r *Runner) Run(phase Phase) error {
	for _, h := range r.hooks {
		if h.Phase != phase {
			continue
		}
		if strings.TrimSpace(h.Command) == "" {
			return errors.New("envhook: hook command must not be empty")
		}
		if err := r.execute(h.Command); err != nil {
			return fmt.Errorf("envhook: hook %q failed: %w", h.Command, err)
		}
	}
	return nil
}

// execute runs a single shell command, wiring its stdout and stderr to the
// runner's configured writers.
func (r *Runner) execute(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = r.stdout
	cmd.Stderr = r.stderr
	return cmd.Run()
}

// Filter returns only the hooks that match the given phase.
func Filter(hooks []Hook, phase Phase) []Hook {
	result := make([]Hook, 0, len(hooks))
	for _, h := range hooks {
		if h.Phase == phase {
			result = append(result, h)
		}
	}
	return result
}
