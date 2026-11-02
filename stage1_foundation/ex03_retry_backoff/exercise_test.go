package ex03_retry_backoff

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

type fakeSleeper struct{ slept []time.Duration }

func (f *fakeSleeper) Sleep(ctx context.Context, d time.Duration) error {
	f.slept = append(f.slept, d)
	return ctx.Err()
}

type scriptedOp struct {
	errs  []error
	calls int
}

func (o *scriptedOp) Run(ctx context.Context) error {
	o.calls++
	if len(o.errs) == 0 {
		return nil
	}
	e := o.errs[0]
	o.errs = o.errs[1:]
	return e
}

type blockingOp struct{}

func (blockingOp) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func TestDelay(t *testing.T) {
	p := RetryPolicy{BaseDelay: 100 * time.Millisecond, MaxDelay: time.Second}
	want := []time.Duration{100, 200, 400, 800, 1000, 1000}
	for i, w := range want {
		if got := p.Delay(i + 1); got != w*time.Millisecond {
			t.Errorf("Delay(%d) = %v, want %v", i+1, got, w*time.Millisecond)
		}
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{nil, false},
		{context.Canceled, false},
		{fmt.Errorf("wrap: %w", context.DeadlineExceeded), false},
		{&HTTPStatusError{Code: 429}, true},
		{&HTTPStatusError{Code: 502}, true},
		{fmt.Errorf("wrap: %w", &HTTPStatusError{Code: 503}), true},
		{&HTTPStatusError{Code: 400}, false},
		{errors.New("connection reset"), true},
	}
	for _, tc := range tests {
		if got := IsRetryable(tc.err); got != tc.want {
			t.Errorf("IsRetryable(%v) = %v, want %v", tc.err, got, tc.want)
		}
	}
}

func TestDo(t *testing.T) {
	p := RetryPolicy{MaxAttempts: 3, BaseDelay: 10 * time.Millisecond, MaxDelay: time.Second}
	e500 := &HTTPStatusError{Code: 500}
	tests := []struct {
		name      string
		errs      []error
		wantCalls int
		wantSleep int
		wantErr   bool
	}{
		{"ok first", nil, 1, 0, false},
		{"ok after two", []error{e500, e500}, 3, 2, false},
		{"exhausted", []error{e500, e500, e500, e500}, 3, 2, true},
		{"non retryable", []error{&HTTPStatusError{Code: 401}}, 1, 0, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			op := &scriptedOp{errs: tc.errs}
			s := &fakeSleeper{}
			err := Do(context.Background(), p, s, op)
			if (err != nil) != tc.wantErr || op.calls != tc.wantCalls || len(s.slept) != tc.wantSleep {
				t.Fatalf("err=%v calls=%d sleeps=%d", err, op.calls, len(s.slept))
			}
			if tc.wantSleep == 2 && (s.slept[0] != 10*time.Millisecond || s.slept[1] != 20*time.Millisecond) {
				t.Fatalf("sleep durations %v", s.slept)
			}
		})
	}
}

func TestDoCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	op := &scriptedOp{errs: []error{errors.New("net"), errors.New("net")}}
	err := Do(ctx, RetryPolicy{MaxAttempts: 5, BaseDelay: time.Millisecond}, &fakeSleeper{}, op)
	if !errors.Is(err, context.Canceled) || op.calls != 1 {
		t.Fatalf("err=%v calls=%d", err, op.calls)
	}
}

func TestCallWithTimeout(t *testing.T) {
	start := time.Now()
	err := CallWithTimeout(context.Background(), 20*time.Millisecond, blockingOp{})
	if !errors.Is(err, context.DeadlineExceeded) || time.Since(start) > time.Second {
		t.Fatalf("err=%v", err)
	}
}
