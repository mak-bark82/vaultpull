package envtag

import (
	"testing"
)

func TestParse_AttachesTags(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "secret", "API_KEY": "abc"}
	annotations := map[string]string{"DB_PASS": "env:prod,tier:db"}
	result := Parse(secrets, annotations)
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	for _, ts := range result {
		if ts.Name == "DB_PASS" {
			if len(ts.Tags) != 2 {
				t.Errorf("expected 2 tags, got %d", len(ts.Tags))
			}
		}
		if ts.Name == "API_KEY" && len(ts.Tags) != 0 {
			t.Errorf("expected no tags for API_KEY")
		}
	}
}

func TestParse_NoAnnotations(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	result := Parse(secrets, map[string]string{})
	if len(result) != 1 || len(result[0].Tags) != 0 {
		t.Errorf("expected 1 untagged secret")
	}
}

func TestFilter_ByTagKey(t *testing.T) {
	secrets := []TaggedSecret{
		{Name: "A", Tags: []Tag{{Key: "env", Value: "prod"}}},
		{Name: "B", Tags: []Tag{{Key: "tier", Value: "db"}}},
		{Name: "C", Tags: []Tag{{Key: "env", Value: "staging"}}},
	}
	result := Filter(secrets, "env")
	if len(result) != 2 {
		t.Errorf("expected 2 filtered secrets, got %d", len(result))
	}
}

func TestFilter_NoMatch(t *testing.T) {
	secrets := []TaggedSecret{
		{Name: "X", Tags: []Tag{{Key: "tier", Value: "web"}}},
	}
	result := Filter(secrets, "env")
	if len(result) != 0 {
		t.Errorf("expected 0 results, got %d", len(result))
	}
}

func TestParseTags_MissingValue(t *testing.T) {
	tags := parseTags("standalone")
	if len(tags) != 1 || tags[0].Key != "standalone" || tags[0].Value != "" {
		t.Errorf("unexpected tag: %+v", tags)
	}
}
