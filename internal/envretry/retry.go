package envretry

import (
	"errors"
	"fmt"
	"time"
)

// Policy defines retry behaviour for vault operations.
type Policy struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64
}

// DefaultPolicy returns a sensible retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		Delay:       200 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// Result holds the outcome of a retried operation.
type Result struct {
	Attempts int
	Err      error
}

// Do executes fn according to p, retrying on transient errors.
// fn should return (true, err) when the error is retryable.
func Do(p Policy, fn func() (bool, error)) Result {
	if p.MaxAttempts < 1 {
		return Result{Attempts: 0, Err: errors.New("envretry: MaxAttempts must be >= 1")}
	}

	delay := p.Delay
	var lastErr error

	for attempt := 1; attempt <= p.MaxAttempts; attempt++ {
		retryable, err := fn()
		if err == nil {
			return Result{Attempts: attempt}
		}
		lastErr = err
		if !retryable || attempt == p.MaxAttempts {
			break
		}
		time.Sleep(delay)
		delay = time.Duration(float64(delay) * p.Multiplier)
	}

	return Result{
		Attempts: p.MaxAttempts,
		Err:      fmt.Errorf("envretry: all %d attempts failed: %w", p.MaxAttempts, lastErr),
	}
}
