package envchain

import (
	"errors"
	"strings"
	"testing"
)

func base() map[string]string {
	return map[string]string{"KEY": "value", "OTHER": "hello"}
}

func TestRun_AllStagesApplied(t *testing.T) {
	c := New()
	c.Add("upper", func(m map[string]string) (map[string]string, error) {
		out := make(map[string]string)
		for k, v := range m {
			out[k] = strings.ToUpper(v)
		}
		return out, nil
	})
	c.Add("prefix", func(m map[string]string) (map[string]string, error) {
		out := make(map[string]string)
		for k, v := range m {
			out[k] = "PRE_" + v
		}
		return out, nil
	})

	results, final, err := c.Run(base())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if final["KEY"] != "PRE_VALUE" {
		t.Errorf("expected PRE_VALUE, got %s", final["KEY"])
	}
}

func TestRun_HaltsOnError(t *testing.T) {
	c := New()
	c.Add("ok", func(m map[string]string) (map[string]string, error) {
		return m, nil
	})
	c.Add("fail", func(m map[string]string) (map[string]string, error) {
		return nil, errors.New("stage failed")
	})
	c.Add("never", func(m map[string]string) (map[string]string, error) {
		t.Error("should not reach this stage")
		return m, nil
	})

	results, _, err := c.Run(base())
	if err == nil {
		t.Fatal("expected error")
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results before halt, got %d", len(results))
	}
}

func TestRun_EmptyChain(t *testing.T) {
	c := New()
	_, final, err := c.Run(base())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if final["KEY"] != "value" {
		t.Errorf("expected original value, got %s", final["KEY"])
	}
}

func TestRun_DoesNotMutateInput(t *testing.T) {
	input := base()
	c := New()
	c.Add("mutate", func(m map[string]string) (map[string]string, error) {
		m["KEY"] = "changed"
		return m, nil
	})
	c.Run(input)
	if input["KEY"] != "value" {
		t.Errorf("input was mutated")
	}
}
