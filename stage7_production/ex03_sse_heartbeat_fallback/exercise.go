//go:build !solution

package ex03_sse_heartbeat_fallback

import (
	"context"
	"io"
	"time"

	"leetagent/common"
)

// WriteSSE 写一个 SSE 事件：
//
//	event 非空时先写 "event: <event>\n"；
//	data 按 "\n" 拆行，每行写 "data: <line>\n"（data 为空也要写一行 "data: \n"）；
//	最后写一个空行 "\n"。
func WriteSSE(w io.Writer, event, data string) error {
	// TODO
	return nil
}

// Heartbeat 定期写 SSE 注释 ": ping\n\n"，防止 nginx / 浏览器 / 负载均衡因空闲断开长连接。
// Run 每从 ticks 收到一次就写一次，直到 ctx 结束或 ticks 关闭；返回写出的心跳次数。写失败立即返回。
// 并发提示：真实场景中心跳和正文写同一个 ResponseWriter，必须共用一把锁——本题通过 Mu 字段体现。
type Heartbeat struct {
	W  io.Writer
	Mu Locker
}

// Locker 就是 sync.Locker，单独声明便于测试注入。
type Locker interface {
	Lock()
	Unlock()
}

func (h *Heartbeat) Run(ctx context.Context, ticks <-chan time.Time) (int, error) {
	// TODO
	return 0, nil
}

// TokenSink 接收流式 token。
type TokenSink interface {
	Token(s string)
}

// Streamer 以流式方式调用某个模型；每收到一段内容调用 sink.Token。
type Streamer interface {
	Stream(ctx context.Context, model string, msgs []common.Message, sink TokenSink) error
}

// StreamWithFallback 先用 primary 模型；如果它失败且**还没有输出过任何 token**，改用 backup 模型重试一次；
// 如果已经输出过 token 才失败，不能切换（用户已经看到半句话），直接返回错误。
// 返回实际使用的模型名。backup 也失败则返回 (backup, err)。
func StreamWithFallback(ctx context.Context, s Streamer, primary, backup string, msgs []common.Message, sink TokenSink) (string, error) {
	// TODO
	return "", nil
}
