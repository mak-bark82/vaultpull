package envfilter

// Filter applies include/exclude key rules to a map of env vars.
// If include is non-empty, only keys in the include list are kept.
// Keys in the exclude list are always removed.
type Filter struct {
	Include []string
	Exclude []string
}

// Apply returns a new map with the filter rules applied.
func (f *Filter) Apply(secrets map[string]string) map[string]string {
	includeSet := toSet(f.Include)
	excludeSet := toSet(f.Exclude)

	result := make(map[string]string)
	for k, v := range secrets {
		if len(includeSet) > 0 {
			if !includeSet[k] {
				continue
			}
		}
		if excludeSet[k] {
			continue
		}
		result[k] = v
	}
	return result
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
