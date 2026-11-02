//go:build !solution

package ex05_slog

import (
	"io"
	"log/slog"
	"time"
)

// RoundStats 是 Agent 循环一轮的统计。
type RoundStats struct {
	RequestID string
	Round     int
	Model     string
	Duration  time.Duration
	Tools     []string // 本轮调用的工具名，按调用顺序
	Err       error
}

// NewJSONLogger 返回写 JSON 到 w 的 logger，最低级别 level。
func NewJSONLogger(w io.Writer, level slog.Level) *slog.Logger {
	// TODO
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// LogRound 记录一轮：
//
//	消息固定为 "agent round"
//	属性：request_id, round, model, duration_ms(int64 毫秒), tool_count(int), tools(逗号拼接的字符串)
//	Err 为 nil 用 Info 级别；否则用 Error 级别并额外带 err 属性（字符串）。
func LogRound(l *slog.Logger, s RoundStats) {
	// TODO
}

// WithRequest 返回一个子 logger，所有日志都自动带上 request_id 和 user 属性。
func WithRequest(l *slog.Logger, requestID, user string) *slog.Logger {
	// TODO
	return l
}
