// Package solution 是 ex05_slog 的参考答案。
package solution

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"time"
)

type RoundStats struct {
	RequestID string
	Round     int
	Model     string
	Duration  time.Duration
	Tools     []string
	Err       error
}

func NewJSONLogger(w io.Writer, level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level}))
}

func LogRound(l *slog.Logger, s RoundStats) {
	attrs := []slog.Attr{
		slog.String("request_id", s.RequestID),
		slog.Int("round", s.Round),
		slog.String("model", s.Model),
		slog.Int64("duration_ms", s.Duration.Milliseconds()),
		slog.Int("tool_count", len(s.Tools)),
		slog.String("tools", strings.Join(s.Tools, ",")),
	}
	level := slog.LevelInfo
	if s.Err != nil {
		level = slog.LevelError
		attrs = append(attrs, slog.String("err", s.Err.Error()))
	}
	l.LogAttrs(context.Background(), level, "agent round", attrs...)
}

func WithRequest(l *slog.Logger, requestID, user string) *slog.Logger {
	return l.With(slog.String("request_id", requestID), slog.String("user", user))
}
