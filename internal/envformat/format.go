package envformat

import (
	"fmt"
	"strings"
)

// Style defines the output format style for env files.
type Style string

const (
	StyleExport  Style = "export"  // export KEY=VALUE
	StylePlain   Style = "plain"   // KEY=VALUE
	StyleQuoted  Style = "quoted"  // KEY="VALUE"
	StyleInline  Style = "inline"  // KEY=VALUE (single line, semicolon-separated)
)

// Options controls formatting behavior.
type Options struct {
	Style     Style
	SortKeys  bool
	Separator string // used for StyleInline; defaults to " ; "
}

// DefaultOptions returns sensible formatting defaults.
func DefaultOptions() Options {
	return Options{
		Style:     StylePlain,
		SortKeys:  true,
		Separator: " ; ",
	}
}

// Format renders a map of env vars into a string using the given options.
func Format(env map[string]string, opts Options) string {
	if len(env) == 0 {
		return ""
	}

	keys := sortedKeys(env, opts.SortKeys)
	lines := make([]string, 0, len(keys))

	for _, k := range keys {
		v := env[k]
		lines = append(lines, renderLine(k, v, opts.Style))
	}

	if opts.Style == StyleInline {
		sep := opts.Separator
		if sep == "" {
			sep = " ; "
		}
		return strings.Join(lines, sep)
	}

	return strings.Join(lines, "\n")
}

func renderLine(key, value string, style Style) string {
	switch style {
	case StyleExport:
		return fmt.Sprintf("export %s=%s", key, value)
	case StyleQuoted:
		return fmt.Sprintf("%s=%q", key, value)
	case StyleInline, StylePlain:
		return fmt.Sprintf("%s=%s", key, value)
	default:
		return fmt.Sprintf("%s=%s", key, value)
	}
}

func sortedKeys(env map[string]string, sort bool) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if sort {
		slices_sort(keys)
	}
	return keys
}

func slices_sort(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
