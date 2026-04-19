package envtag

import (
	"strings"
)

// Tag represents metadata attached to a secret key.
type Tag struct {
	Key   string
	Value string
}

// TaggedSecret holds a secret key/value with associated tags.
type TaggedSecret struct {
	Name  string
	Value string
	Tags  []Tag
}

// Parse reads a map of secrets and a map of tag annotations,
// returning a slice of TaggedSecrets.
// Tag annotations are expected in the form: "KEY=tag1:val1,tag2:val2"
func Parse(secrets map[string]string, annotations map[string]string) []TaggedSecret {
	result := make([]TaggedSecret, 0, len(secrets))
	for name, value := range secrets {
		ts := TaggedSecret{Name: name, Value: value}
		if raw, ok := annotations[name]; ok {
			ts.Tags = parseTags(raw)
		}
		result = append(result, ts)
	}
	return result
}

// Filter returns only those TaggedSecrets that have the given tag key.
func Filter(secrets []TaggedSecret, tagKey string) []TaggedSecret {
	out := []TaggedSecret{}
	for _, s := range secrets {
		for _, t := range s.Tags {
			if t.Key == tagKey {
				out = append(out, s)
				break
			}
		}
	}
	return out
}

func parseTags(raw string) []Tag {
	tags := []Tag{}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, ":", 2)
		t := Tag{Key: strings.TrimSpace(kv[0])}
		if len(kv) == 2 {
			t.Value = strings.TrimSpace(kv[1])
		}
		tags = append(tags, t)
	}
	return tags
}
