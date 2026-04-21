package envpatch

import "fmt"

// Op represents the type of patch operation.
type Op string

const (
	OpSet    Op = "set"
	OpDelete Op = "delete"
	OpRename Op = "rename"
)

// Patch describes a single mutation to apply to a secret map.
type Patch struct {
	Op      Op     `yaml:"op"`
	Key     string `yaml:"key"`
	Value   string `yaml:"value,omitempty"`
	NewKey  string `yaml:"new_key,omitempty"`
}

// Result holds the outcome of applying a patch.
type Result struct {
	Applied []string
	Skipped []string
}

// Apply applies a slice of Patch operations to the provided secrets map.
// It returns a new map with all mutations applied and a Result summary.
func Apply(secrets map[string]string, patches []Patch) (map[string]string, Result, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var result Result

	for _, p := range patches {
		if p.Key == "" {
			return nil, Result{}, fmt.Errorf("patch missing required field 'key'")
		}

		switch p.Op {
		case OpSet:
			out[p.Key] = p.Value
			result.Applied = append(result.Applied, p.Key)

		case OpDelete:
			if _, exists := out[p.Key]; !exists {
				result.Skipped = append(result.Skipped, p.Key)
				continue
			}
			delete(out, p.Key)
			result.Applied = append(result.Applied, p.Key)

		case OpRename:
			if p.NewKey == "" {
				return nil, Result{}, fmt.Errorf("rename patch for key %q missing 'new_key'", p.Key)
			}
			val, exists := out[p.Key]
			if !exists {
				result.Skipped = append(result.Skipped, p.Key)
				continue
			}
			delete(out, p.Key)
			out[p.NewKey] = val
			result.Applied = append(result.Applied, p.Key)

		default:
			return nil, Result{}, fmt.Errorf("unknown patch op: %q", p.Op)
		}
	}

	return out, result, nil
}
