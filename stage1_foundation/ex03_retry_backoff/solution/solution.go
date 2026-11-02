// Package solution 是 ex03_retry_backoff 的参考答案。
package solution

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type HTTPStatusError struct{ Code int }

func (e *HTTPStatusError) Error() string { return fmt.Sprintf("http status %d", e.Code) }

type Operation interface {
	Run(ctx context.Context) error
}

type Sleeper interface {
	Sleep(ctx context.Context, d time.Duration) error
}

type RetryPolicy struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

func (p RetryPolicy) Delay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	d := p.BaseDelay
	for i := 1; i < attempt; i++ {
		d *= 2
		if p.MaxDelay > 0 && d >= p.MaxDelay {
			return p.MaxDelay
		}
	}
	if p.MaxDelay > 0 && d > p.MaxDelay {
		return p.MaxDelay
	}
	return d
}

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	var se *HTTPStatusError
	if errors.As(err, &se) {
		return se.Code == 429 || se.Code >= 500
	}
	return true
}

func Do(ctx context.Context, p RetryPolicy, s Sleeper, op Operation) error {
	attempts := p.MaxAttempts
	if attempts <= 0 {
		attempts = 1
	}
	var err error
	for i := 1; i <= attempts; i++ {
		err = op.Run(ctx)
		if err == nil || !IsRetryable(err) {
			return err
		}
		if i == attempts {
			break
		}
		if serr := s.Sleep(ctx, p.Delay(i)); serr != nil {
			return serr
		}
	}
	return err
}

func CallWithTimeout(ctx context.Context, d time.Duration, op Operation) error {
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()
	return op.Run(ctx)
}
