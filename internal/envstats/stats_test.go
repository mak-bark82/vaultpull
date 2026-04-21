package envstats_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/envstats"
)

func TestCompute_EmptyMap(t *testing.T) {
	s := envstats.Compute(map[string]string{})
	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
}

func TestCompute_NilMap(t *testing.T) {
	s := envstats.Compute(nil)
	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
}

func TestCompute_Counts(t *testing.T) {
	env := map[string]string{
		"KEY_A": "hello",
		"KEY_B": "",
		"KEY_C": "world!",
	}
	s := envstats.Compute(env)
	if s.Total != 3 {
		t.Errorf("expected Total=3, got %d", s.Total)
	}
	if s.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", s.Empty)
	}
	if s.NonEmpty != 2 {
		t.Errorf("expected NonEmpty=2, got %d", s.NonEmpty)
	}
}

func TestCompute_Lengths(t *testing.T) {
	env := map[string]string{
		"SHORT": "hi",
		"LONG":  "averylongvalue",
		"MID":   "middle",
	}
	s := envstats.Compute(env)
	if s.MinLength != 2 {
		t.Errorf("expected MinLength=2, got %d", s.MinLength)
	}
	if s.MaxLength != 14 {
		t.Errorf("expected MaxLength=14, got %d", s.MaxLength)
	}
	expectedAvg := float64(2+14+6) / 3.0
	if s.AvgLength != expectedAvg {
		t.Errorf("expected AvgLength=%.4f, got %.4f", expectedAvg, s.AvgLength)
	}
}

func TestCompute_LongestAndShortestKey(t *testing.T) {
	env := map[string]string{
		"A": "short",
		"B": "this is a much longer value here",
	}
	s := envstats.Compute(env)
	if s.ShortestKey != "A" {
		t.Errorf("expected ShortestKey=A, got %s", s.ShortestKey)
	}
	if s.LongestKey != "B" {
		t.Errorf("expected LongestKey=B, got %s", s.LongestKey)
	}
}

func TestCompute_SingleEntry(t *testing.T) {
	env := map[string]string{"ONLY": "value"}
	s := envstats.Compute(env)
	if s.Total != 1 {
		t.Errorf("expected Total=1, got %d", s.Total)
	}
	if s.MinLength != s.MaxLength {
		t.Errorf("expected MinLength==MaxLength for single entry")
	}
	if s.LongestKey != s.ShortestKey {
		t.Errorf("expected LongestKey==ShortestKey for single entry")
	}
}
