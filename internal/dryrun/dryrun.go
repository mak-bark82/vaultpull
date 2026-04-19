package dryrun

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultpull/internal/diff"
)

// Reporter prints a dry-run diff report without writing any files.
type Reporter struct {
	out io.Writer
}

// NewReporter returns a Reporter that writes to out.
// If out is nil, os.Stdout is used.
func NewReporter(out io.Writer) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{out: out}
}

// Report prints a human-readable diff between existing and incoming secrets
// for the given env file path.
func (r *Reporter) Report(envFile string, existing, incoming map[string]string) {
	changes := diff.Compare(existing, incoming)
	if len(changes) == 0 {
		fmt.Fprintf(r.out, "[dry-run] %s: no changes\n", envFile)
		return
	}
	fmt.Fprintf(r.out, "[dry-run] %s:\n", envFile)
	for _, c := range changes {
		switch c.Status {
		case diff.Added:
			fmt.Fprintf(r.out, "  + %s\n", c.Key)
		case diff.Removed:
			fmt.Fprintf(r.out, "  - %s\n", c.Key)
		case diff.Changed:
			fmt.Fprintf(r.out, "  ~ %s\n", c.Key)
		}
	}
	fmt.Fprintf(r.out, "  %s\n", diff.Summary(changes))
}
