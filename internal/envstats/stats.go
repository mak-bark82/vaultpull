package envstats

import "sort"

// Stats holds aggregate statistics about a set of environment secrets.
type Stats struct {
	Total     int
	Empty     int
	NonEmpty  int
	AvgLength float64
	MinLength int
	MaxLength int
	LongestKey string
	ShortestKey string
}

// Compute calculates statistics over the provided key-value map.
// Returns a zero-value Stats if the map is nil or empty.
func Compute(env map[string]string) Stats {
	if len(env) == 0 {
		return Stats{}
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var totalLen int
	minLen := len(env[keys[0]])
	maxLen := len(env[keys[0]])
	minKey := keys[0]
	maxKey := keys[0]
	empty := 0

	for _, k := range keys {
		v := env[k]
		vl := len(v)
		if vl == 0 {
			empty++
		}
		totalLen += vl
		if vl < minLen {
			minLen = vl
			minKey = k
		}
		if vl > maxLen {
			maxLen = vl
			maxKey = k
		}
	}

	return Stats{
		Total:       len(keys),
		Empty:       empty,
		NonEmpty:    len(keys) - empty,
		AvgLength:   float64(totalLen) / float64(len(keys)),
		MinLength:   minLen,
		MaxLength:   maxLen,
		ShortestKey: minKey,
		LongestKey:  maxKey,
	}
}
