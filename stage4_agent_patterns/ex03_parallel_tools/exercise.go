//go:build !solution

package ex03_parallel_tools

import (
	"context"

	"leetagent/common"
)

// Result 是一个工具调用的执行结果。
type Result struct {
	CallID  string
	Content string // 出错时为 "error: " + err.Error()
	Err     error
}

// ExecuteAll 并发执行一批工具调用（模型一次返回多个 tool_calls 时很常见）。
//   - 同时运行的调用数不超过 maxConcurrency（<=0 视为 1）；
//   - 返回切片顺序必须与 calls 一致（tool 消息必须按 tool_calls 的顺序回填）；
//   - 单个调用出错不影响其他调用；
//   - ctx 被取消后，尚未开始的调用不再执行，其 Err 为 ctx.Err()，Content 为 "error: " + ctx.Err().Error()。
//
// 约束：不要使用闭包，worker 用具名函数/方法实现。
func ExecuteAll(ctx context.Context, runner common.ToolRunner, calls []common.ToolCall, maxConcurrency int) []Result {
	// TODO
	return nil
}

// ToMessages 把结果转成 role=tool 的消息。
func ToMessages(results []Result) []common.Message {
	// TODO
	return nil
}
