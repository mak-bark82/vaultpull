// Package envcast provides utilities for casting environment variable
// string values into typed Go values with error reporting.
package envcast

import (
	"fmt"
	"strconv"
	"strings"
)

// Result holds the outcome of a cast operation for a single key.
type Result struct {
	Key    string
	Raw    string
	Typed  interface{}
	Err    error
}

// Options controls cast behaviour.
type Options struct {
	// Types maps key names to target type strings: "bool", "int", "float", "string".
	Types map[string]string
}

// Cast iterates over secrets and casts each value according to the type
// mapping in opts. Keys absent from the mapping are returned as-is (string).
// All results are returned; callers should inspect Result.Err for failures.
func Cast(secrets map[string]string, opts Options) []Result {
	results := make([]Result, 0, len(secrets))
	for k, v := range secrets {
		r := Result{Key: k, Raw: v}
		target, ok := opts.Types[k]
		if !ok {
			r.Typed = v
			results = append(results, r)
			continue
		}
		typed, err := castValue(v, strings.ToLower(strings.TrimSpace(target)))
		r.Typed = typed
		r.Err = err
		results = append(results, r)
	}
	return results
}

// CastOne casts a single value to the named type.
func CastOne(key, value, typeName string) Result {
	r := Result{Key: key, Raw: value}
	typed, err := castValue(value, strings.ToLower(strings.TrimSpace(typeName)))
	r.Typed = typed
	r.Err = err
	return r
}

func castValue(v, typeName string) (interface{}, error) {
	switch typeName {
	case "bool":
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("cannot cast %q to bool: %w", v, err)
		}
		return b, nil
	case "int":
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot cast %q to int: %w", v, err)
		}
		return i, nil
	case "float":
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot cast %q to float: %w", v, err)
		}
		return f, nil
	case "string":
		return v, nil
	default:
		return nil, fmt.Errorf("unknown type %q", typeName)
	}
}
