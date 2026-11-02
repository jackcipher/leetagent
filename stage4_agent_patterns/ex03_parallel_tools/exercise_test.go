package ex03_parallel_tools

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"leetagent/common"
)

type slowRunner struct {
	mu      sync.Mutex
	running int32
	peak    int32
	calls   int32
	delay   time.Duration
	onStart func(name string)
}

func (s *slowRunner) Call(ctx context.Context, name, args string) (string, error) {
	atomic.AddInt32(&s.calls, 1)
	n := atomic.AddInt32(&s.running, 1)
	s.mu.Lock()
	if n > s.peak {
		s.peak = n
	}
	s.mu.Unlock()
	if s.onStart != nil {
		s.onStart(name)
	}
	time.Sleep(s.delay)
	atomic.AddInt32(&s.running, -1)
	if name == "bad" {
		return "", errors.New("boom")
	}
	return "ok:" + name + args, nil
}

func makeCalls(n int) []common.ToolCall {
	var calls []common.ToolCall
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("t%d", i)
		if i == 2 {
			name = "bad"
		}
		calls = append(calls, common.Call(fmt.Sprintf("c%d", i), name, "{}"))
	}
	return calls
}

func TestOrderAndErrors(t *testing.T) {
	r := &slowRunner{delay: 5 * time.Millisecond}
	res := ExecuteAll(context.Background(), r, makeCalls(6), 3)
	if len(res) != 6 {
		t.Fatalf("len = %d", len(res))
	}
	for i, x := range res {
		if x.CallID != fmt.Sprintf("c%d", i) {
			t.Fatalf("order broken at %d: %+v", i, x)
		}
	}
	if res[2].Err == nil || res[2].Content != "error: boom" || res[0].Content != "ok:t0{}" {
		t.Fatalf("results = %+v", res)
	}
	msgs := ToMessages(res)
	if msgs[2].Role != "tool" || msgs[2].ToolCallID != "c2" || msgs[2].Content != "error: boom" {
		t.Fatalf("msgs = %+v", msgs)
	}
}

func TestConcurrencyLimit(t *testing.T) {
	tests := []struct{ limit, wantPeak int32 }{{1, 1}, {3, 3}, {0, 1}}
	for _, tc := range tests {
		r := &slowRunner{delay: 20 * time.Millisecond}
		start := time.Now()
		ExecuteAll(context.Background(), r, makeCalls(6), int(tc.limit))
		if r.peak != tc.wantPeak {
			t.Errorf("limit %d: peak = %d", tc.limit, r.peak)
		}
		if tc.limit == 3 && time.Since(start) > 100*time.Millisecond {
			t.Errorf("limit 3 should run in ~40ms, took %v", time.Since(start))
		}
	}
}

func TestCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	r := &slowRunner{delay: 10 * time.Millisecond}
	r.onStart = func(string) { cancel() }
	res := ExecuteAll(ctx, r, makeCalls(5), 1)
	if r.calls != 1 {
		t.Fatalf("only the first call should run, got %d", r.calls)
	}
	for _, x := range res[1:] {
		if !errors.Is(x.Err, context.Canceled) || x.Content != "error: context canceled" {
			t.Fatalf("res = %+v", x)
		}
	}
}
