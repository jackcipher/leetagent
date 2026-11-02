//go:build !solution

package ex04_approval_gate

import (
	"context"
	"errors"

	"leetagent/common"
)

var ErrDenied = errors.New("tool call denied by human")

// Approver 决定是否允许一个危险调用。
type Approver interface {
	Approve(ctx context.Context, call common.ToolCall) (bool, error)
}

// Gate 包装一个 ToolRunner：Dangerous 中的工具必须经过 Approver 同意才执行（Human-in-the-loop）。
type Gate struct {
	Inner     common.ToolRunner
	Approver  Approver
	Dangerous map[string]bool
	Log       []string // 每次决策追加 "<name>: approved" / "<name>: denied" / "<name>: auto"
}

// Call 实现 common.ToolRunner。
//
//	非危险工具：直接执行，Log 追加 "<name>: auto"。
//	危险工具：调用 Approver（call 的 ID 留空，Function 填 name/args）；
//	  同意 -> 执行，Log "<name>: approved"；拒绝 -> 返回 ("", ErrDenied)，Log "<name>: denied"；
//	  Approver 出错 -> 原样返回该错误，不执行、不记录。
func (g *Gate) Call(ctx context.Context, name, args string) (string, error) {
	// TODO
	return "", nil
}

// Pending 是一个等待人类决定的请求。
type Pending struct {
	ID   string
	Call common.ToolCall
}

// ChanApprover 通过 channel 把审批请求交给前端/人类（对应 SSE 推一个 approval 事件，然后等 POST /approve）。
//
//	NewChanApprover 创建；Requests() 返回只读 channel，每个请求的 ID 依次为 "req-1"、"req-2"...
//	Approve 发出请求后阻塞，直到 Resolve(id, ok) 被调用或 ctx 结束（返回 false, ctx.Err()）。
//	Resolve 对未知 id 返回 false。
type ChanApprover struct {
	// TODO: 字段自行设计（需要互斥锁）
}

func NewChanApprover() *ChanApprover {
	// TODO
	return &ChanApprover{}
}

func (a *ChanApprover) Requests() <-chan Pending {
	// TODO
	return nil
}

func (a *ChanApprover) Approve(ctx context.Context, call common.ToolCall) (bool, error) {
	// TODO
	return false, nil
}

func (a *ChanApprover) Resolve(id string, ok bool) bool {
	// TODO
	return false
}
