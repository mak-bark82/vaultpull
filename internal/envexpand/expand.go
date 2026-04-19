package envexpand

import (
	"os"
	"strings"
)

// Expand resolves variable references within env values.
// Values may reference other keys in the same map using ${VAR} or $VAR syntax.
// References not found in the map fall back to os.Getenv.
func Expand(env map[string]string) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = expandValue(v, env)
	}
	return result
}

// ExpandValue expands a single value string against the provided env map.
func ExpandValue(value string, env map[string]string) string {
	return expandValue(value, env)
}

func expandValue(value string, env map[string]string) string {
	return os.Expand(value, func(key string) string {
		if val, ok := env[key]; ok {
			return val
		}
		return os.Getenv(key)
	})
}

// HasReferences returns true if the value contains variable references.
func HasReferences(value string) bool {
	return strings.Contains(value, "$")
}
