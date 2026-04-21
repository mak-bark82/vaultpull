package envcompare_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envcompare"
)

func TestCompare_OnlyInLeft(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"B": "2"}
	r := envcompare.Compare(left, right)
	if _, ok := r.OnlyInLeft["A"]; !ok {
		t.Error("expected A to be only-left")
	}
	if len(r.OnlyInRight) != 0 {
		t.Errorf("unexpected only-right keys: %v", r.OnlyInRight)
	}
}

func TestCompare_OnlyInRight(t *testing.T) {
	left := map[string]string{"B": "2"}
	right := map[string]string{"B": "2", "C": "3"}
	r := envcompare.Compare(left, right)
	if _, ok := r.OnlyInRight["C"]; !ok {
		t.Error("expected C to be only-right")
	}
}

func TestCompare_Different(t *testing.T) {
	left := map[string]string{"X": "old"}
	right := map[string]string{"X": "new"}
	r := envcompare.Compare(left, right)
	pair, ok := r.Different["X"]
	if !ok {
		t.Fatal("expected X in Different")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("unexpected pair: %v", pair)
	}
}

func TestCompare_Identical(t *testing.T) {
	left := map[string]string{"K": "v"}
	right := map[string]string{"K": "v"}
	r := envcompare.Compare(left, right)
	if _, ok := r.Identical["K"]; !ok {
		t.Error("expected K in Identical")
	}
	if len(r.Different) != 0 || len(r.OnlyInLeft) != 0 || len(r.OnlyInRight) != 0 {
		t.Error("unexpected non-identical results")
	}
}

func TestRender_ContainsSymbols(t *testing.T) {
	left := map[string]string{"A": "1", "B": "old"}
	right := map[string]string{"B": "new", "C": "3"}
	r := envcompare.Compare(left, right)
	var buf bytes.Buffer
	envcompare.Render(&buf, r, false)
	out := buf.String()
	if !strings.Contains(out, "< A=") {
		t.Errorf("expected '< A=' in output, got: %s", out)
	}
	if !strings.Contains(out, "> C=") {
		t.Errorf("expected '> C=' in output, got: %s", out)
	}
	if !strings.Contains(out, "~ B:") {
		t.Errorf("expected '~ B:' in output, got: %s", out)
	}
}

func TestRender_Redact(t *testing.T) {
	left := map[string]string{"SECRET": "mysecret"}
	right := map[string]string{"SECRET": "other"}
	r := envcompare.Compare(left, right)
	var buf bytes.Buffer
	envcompare.Render(&buf, r, true)
	if strings.Contains(buf.String(), "mysecret") {
		t.Error("redacted render should not contain plaintext value")
	}
}

func TestSummary_NoDifferences(t *testing.T) {
	r := envcompare.Compare(map[string]string{"K": "v"}, map[string]string{"K": "v"})
	if s := envcompare.Summary(r); s != "1 identical" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_Mixed(t *testing.T) {
	left := map[string]string{"A": "1", "B": "old"}
	right := map[string]string{"B": "new", "C": "3"}
	r := envcompare.Compare(left, right)
	s := envcompare.Summary(r)
	if !strings.Contains(s, "only-left") || !strings.Contains(s, "only-right") || !strings.Contains(s, "changed") {
		t.Errorf("unexpected summary: %s", s)
	}
}
