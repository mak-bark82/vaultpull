package envcompare

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Result holds the outcome of comparing two env maps.
type Result struct {
	OnlyInLeft  map[string]string
	OnlyInRight map[string]string
	Different   map[string][2]string // key -> [left, right]
	Identical   map[string]string
}

// Compare performs a full key-by-key comparison of two env maps.
func Compare(left, right map[string]string) Result {
	r := Result{
		OnlyInLeft:  make(map[string]string),
		OnlyInRight: make(map[string]string),
		Different:   make(map[string][2]string),
		Identical:   make(map[string]string),
	}

	for k, lv := range left {
		if rv, ok := right[k]; ok {
			if lv == rv {
				r.Identical[k] = lv
			} else {
				r.Different[k] = [2]string{lv, rv}
			}
		} else {
			r.OnlyInLeft[k] = lv
		}
	}

	for k, rv := range right {
		if _, ok := left[k]; !ok {
			r.OnlyInRight[k] = rv
		}
	}

	return r
}

// Render writes a human-readable comparison report to w.
func Render(w io.Writer, r Result, redact bool) {
	mask := func(v string) string {
		if redact {
			return "***"
		}
		return v
	}

	keys := func(m map[string]string) []string {
		out := make([]string, 0, len(m))
		for k := range m {
			out = append(out, k)
		}
		sort.Strings(out)
		return out
	}

	for _, k := range keys(r.OnlyInLeft) {
		fmt.Fprintf(w, "< %s=%s\n", k, mask(r.OnlyInLeft[k]))
	}
	for _, k := range keys(r.OnlyInRight) {
		fmt.Fprintf(w, "> %s=%s\n", k, mask(r.OnlyInRight[k]))
	}

	diffKeys := make([]string, 0, len(r.Different))
	for k := range r.Different {
		diffKeys = append(diffKeys, k)
	}
	sort.Strings(diffKeys)
	for _, k := range diffKeys {
		pair := r.Different[k]
		fmt.Fprintf(w, "~ %s: %s -> %s\n", k, mask(pair[0]), mask(pair[1]))
	}
}

// Summary returns a one-line summary string for the result.
func Summary(r Result) string {
	parts := []string{}
	if n := len(r.OnlyInLeft); n > 0 {
		parts = append(parts, fmt.Sprintf("%d only-left", n))
	}
	if n := len(r.OnlyInRight); n > 0 {
		parts = append(parts, fmt.Sprintf("%d only-right", n))
	}
	if n := len(r.Different); n > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", n))
	}
	if n := len(r.Identical); n > 0 {
		parts = append(parts, fmt.Sprintf("%d identical", n))
	}
	if len(parts) == 0 {
		return "no differences"
	}
	return strings.Join(parts, ", ")
}
