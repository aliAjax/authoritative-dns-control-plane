package operations

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type RetryPolicy struct {
	Attempts int
	Initial  time.Duration
	Max      time.Duration
	Jitter   float64
}

func DefaultRetry() RetryPolicy {
	return RetryPolicy{Attempts: 5, Initial: 50 * time.Millisecond, Max: 2 * time.Second, Jitter: .2}
}
func Do(ctx context.Context, p RetryPolicy, fn func(context.Context) error) error {
	if p.Attempts < 1 {
		p.Attempts = 1
	}
	if p.Initial <= 0 {
		p.Initial = time.Millisecond
	}
	if p.Max <= 0 {
		p.Max = time.Second
	}
	var last error
	delay := p.Initial
	for i := 0; i < p.Attempts; i++ {
		if err := fn(ctx); err == nil {
			return nil
		} else {
			last = err
		}
		if i == p.Attempts-1 {
			break
		}
		j := 1 + (rand.Float64()*2-1)*p.Jitter
		wait := time.Duration(float64(delay) * j)
		tm := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			tm.Stop()
			return ctx.Err()
		case <-tm.C:
		}
		delay *= 2
		if delay > p.Max {
			delay = p.Max
		}
	}
	return fmt.Errorf("retry exhausted: %w", last)
}
func Backoff(p RetryPolicy, attempt int) time.Duration {
	d := p.Initial
	for i := 0; i < attempt; i++ {
		d *= 2
		if d >= p.Max {
			return p.Max
		}
	}
	return d
}
