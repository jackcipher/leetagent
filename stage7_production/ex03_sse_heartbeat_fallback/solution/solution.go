// Package solution 是 ex03_sse_heartbeat_fallback 的参考答案。
package solution

import (
	"context"
	"io"
	"strings"
	"time"

	"leetagent/common"
)

func WriteSSE(w io.Writer, event, data string) error {
	var b strings.Builder
	if event != "" {
		b.WriteString("event: " + event + "\n")
	}
	for _, line := range strings.Split(data, "\n") {
		b.WriteString("data: " + line + "\n")
	}
	b.WriteString("\n")
	_, err := io.WriteString(w, b.String())
	return err
}

type Locker interface {
	Lock()
	Unlock()
}

type Heartbeat struct {
	W  io.Writer
	Mu Locker
}

func (h *Heartbeat) ping() error {
	h.Mu.Lock()
	defer h.Mu.Unlock()
	_, err := io.WriteString(h.W, ": ping\n\n")
	return err
}

func (h *Heartbeat) Run(ctx context.Context, ticks <-chan time.Time) (int, error) {
	n := 0
	for {
		select {
		case <-ctx.Done():
			return n, nil
		case _, ok := <-ticks:
			if !ok {
				return n, nil
			}
			if err := h.ping(); err != nil {
				return n, err
			}
			n++
		}
	}
}

type TokenSink interface {
	Token(s string)
}

type Streamer interface {
	Stream(ctx context.Context, model string, msgs []common.Message, sink TokenSink) error
}

type countingSink struct {
	inner TokenSink
	n     int
}

func (c *countingSink) Token(s string) {
	c.n++
	c.inner.Token(s)
}

func StreamWithFallback(ctx context.Context, s Streamer, primary, backup string, msgs []common.Message, sink TokenSink) (string, error) {
	cs := &countingSink{inner: sink}
	err := s.Stream(ctx, primary, msgs, cs)
	if err == nil {
		return primary, nil
	}
	if cs.n > 0 || ctx.Err() != nil {
		return primary, err
	}
	return backup, s.Stream(ctx, backup, msgs, sink)
}
