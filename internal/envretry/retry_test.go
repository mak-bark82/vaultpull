package envretry_test

import (
	"errors"
	"testing"
	"time"

	"github.com/user/vaultpull/internal/envretry"
)

func TestDo_SucceedsFirstAttempt(t *testing.T) {
	calls := 0
	r := envretry.Do(envretry.DefaultPolicy(), func() (bool, error) {
		calls++
		return false, nil
	})
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
	if r.Attempts != 1 {
		t.Fatalf("expected Attempts=1, got %d", r.Attempts)
	}
}

func TestDo_RetriesOnTransientError(t *testing.T) {
	p := envretry.Policy{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}
	calls := 0
	transient := errors.New("temporary failure")

	r := envretry.Do(p, func() (bool, error) {
		calls++
		if calls < 3 {
			return true, transient
		}
		return false, nil
	})
	if r.Err != nil {
		t.Fatalf("unexpected error: %v", r.Err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_StopsOnNonRetryableError(t *testing.T) {
	p := envretry.Policy{MaxAttempts: 5, Delay: time.Millisecond, Multiplier: 1.0}
	calls := 0

	r := envretry.Do(p, func() (bool, error) {
		calls++
		return false, errors.New("fatal error")
	})
	if r.Err == nil {
		t.Fatal("expected error, got nil")
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_ExhaustsAllAttempts(t *testing.T) {
	p := envretry.Policy{MaxAttempts: 4, Delay: time.Millisecond, Multiplier: 1.0}
	calls := 0

	r := envretry.Do(p, func() (bool, error) {
		calls++
		return true, errors.New("always fails")
	})
	if r.Err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
	if calls != 4 {
		t.Fatalf("expected 4 calls, got %d", calls)
	}
}

func TestDo_InvalidMaxAttempts(t *testing.T) {
	p := envretry.Policy{MaxAttempts: 0, Delay: time.Millisecond, Multiplier: 1.0}
	r := envretry.Do(p, func() (bool, error) { return false, nil })
	if r.Err == nil {
		t.Fatal("expected error for MaxAttempts=0")
	}
}

func TestDefaultPolicy_Values(t *testing.T) {
	p := envretry.DefaultPolicy()
	if p.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", p.MaxAttempts)
	}
	if p.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", p.Multiplier)
	}
}
