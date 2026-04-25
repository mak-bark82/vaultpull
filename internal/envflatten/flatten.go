package envflatten

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls how nested keys are flattened.
type Options struct {
	// Separator is placed between key segments (default: "_").
	Separator string
	// UpperCase converts all resulting keys to uppercase.
	UpperCase bool
	// Prefix is prepended to every flattened key.
	Prefix string
}

// DefaultOptions returns sensible defaults for flattening.
func DefaultOptions() Options {
	return Options{
		Separator: "_",
		UpperCase: true,
	}
}

// Flatten takes a nested map (map values may themselves be
// map[string]interface{} or scalar types) and returns a flat
// map[string]string suitable for use as env vars.
func Flatten(input map[string]interface{}, opts Options) (map[string]string, error) {
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	result := make(map[string]string)
	if err := flatten(input, opts.Prefix, opts.Separator, result); err != nil {
		return nil, err
	}
	if opts.UpperCase {
		upped := make(map[string]string, len(result))
		for k, v := range result {
			upped[strings.ToUpper(k)] = v
		}
		return upped, nil
	}
	return result, nil
}

func flatten(input map[string]interface{}, prefix, sep string, out map[string]string) error {
	keys := make([]string, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := input[k]
		fullKey := k
		if prefix != "" {
			fullKey = prefix + sep + k
		}
		switch val := v.(type) {
		case map[string]interface{}:
			if err := flatten(val, fullKey, sep, out); err != nil {
				return err
			}
		case string:
			out[fullKey] = val
		case nil:
			out[fullKey] = ""
		default:
			out[fullKey] = fmt.Sprintf("%v", val)
		}
	}
	return nil
}
