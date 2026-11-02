package ex03_sse_heartbeat_fallback

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"leetagent/common"
)

func TestWriteSSE(t *testing.T) {
	tests := []struct{ event, data, want string }{
		{"", `{"delta":"你好"}`, "data: {\"delta\":\"你好\"}\n\n"},
		{"tool", "a\nb", "event: tool\ndata: a\ndata: b\n\n"},
		{"done", "", "event: done\ndata: \n\n"},
	}
	for _, tc := range tests {
		var b strings.Builder
		if err := WriteSSE(&b, tc.event, tc.data); err != nil || b.String() != tc.want {
			t.Errorf("got %q err=%v want %q", b.String(), err, tc.want)
		}
	}
}

type countingLock struct {
	sync.Mutex
	n int
}

func (c *countingLock) Lock() { c.Mutex.Lock(); c.n++ }

func TestHeartbeat(t *testing.T) {
	var b strings.Builder
	lock := &countingLock{}
	h := &Heartbeat{W: &b, Mu: lock}
	ticks := make(chan time.Time, 3)
	for i := 0; i < 3; i++ {
		ticks <- time.Time{}
	}
	close(ticks)
	n, err := h.Run(context.Background(), ticks)
	if err != nil || n != 3 || b.String() != strings.Repeat(": ping\n\n", 3) {
		t.Fatalf("n=%d err=%v out=%q", n, err, b.String())
	}
	if lock.n != 3 {
		t.Fatalf("each ping must take the lock, got %d", lock.n)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	never := make(chan time.Time)
	if n, _ := h.Run(ctx, never); n != 0 {
		t.Fatalf("cancelled ctx: n=%d", n)
	}
}

type sinkRec struct{ got []string }

func (s *sinkRec) Token(x string) { s.got = append(s.got, x) }

type scriptedStreamer struct {
	plan  map[string][]string // model -> tokens
	fails map[string]bool     // model -> fail after tokens
	calls []string
}

func (s *scriptedStreamer) Stream(ctx context.Context, model string, msgs []common.Message, sink TokenSink) error {
	s.calls = append(s.calls, model)
	for _, tok := range s.plan[model] {
		sink.Token(tok)
	}
	if s.fails[model] {
		return errors.New(model + " failed")
	}
	return nil
}

func TestStreamWithFallback(t *testing.T) {
	tests := []struct {
		name      string
		plan      map[string][]string
		fails     map[string]bool
		wantModel string
		wantErr   string
		wantToks  []string
		wantCalls []string
	}{
		{"primary ok", map[string][]string{"big": {"你", "好"}}, nil, "big", "", []string{"你", "好"}, []string{"big"}},
		{"fail before first token", map[string][]string{"small": {"嗨"}}, map[string]bool{"big": true}, "small", "", []string{"嗨"}, []string{"big", "small"}},
		{"fail mid-stream", map[string][]string{"big": {"半句"}}, map[string]bool{"big": true}, "big", "big failed", []string{"半句"}, []string{"big"}},
		{"both fail", nil, map[string]bool{"big": true, "small": true}, "small", "small failed", nil, []string{"big", "small"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &scriptedStreamer{plan: tc.plan, fails: tc.fails}
			sink := &sinkRec{}
			model, err := StreamWithFallback(context.Background(), s, "big", "small", nil, sink)
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			if model != tc.wantModel || gotErr != tc.wantErr || !reflect.DeepEqual(sink.got, tc.wantToks) || !reflect.DeepEqual(s.calls, tc.wantCalls) {
				t.Fatalf("model=%q err=%q toks=%v calls=%v", model, gotErr, sink.got, s.calls)
			}
		})
	}
}
