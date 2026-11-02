//go:build !solution

package ex03_subagent

import (
	"context"
	"errors"

	"leetagent/common"
)

var (
	ErrMaxRounds   = errors.New("subagent exceeded max rounds")
	ErrBadDelegate = errors.New("delegate arguments must be {\"task\": \"...\"} with non-empty task")
)

// SubAgent 是一个专精的子 Agent：自己的 system prompt、自己的工具集，独立上下文。
type SubAgent struct {
	Name        string // 例如 "researcher"
	Description string // 给主 Agent 看的能力说明
	System      string
	LLM         common.Chatter
	Model       string
	Tools       common.ToolRunner // 可以为 nil（没有工具）
	Defs        []common.Tool
	MaxRounds   int
}

// Run 用全新上下文 [{system: System}, {user: task}] 执行 ReAct 循环（同 stage4 ex01：工具错误回填 "error: ..."），
// 返回最终文本。超出 MaxRounds 返回 ErrMaxRounds。
func (s *SubAgent) Run(ctx context.Context, task string) (string, error) {
	// TODO
	return "", nil
}

// DelegateTool 生成主 Agent 用来委派任务的工具定义：
//
//	name = "delegate_" + s.Name
//	description = s.Description
//	parameters = {"type":"object","properties":{"task":{"type":"string","description":"交给子 Agent 的完整任务描述"}},"required":["task"]}
func DelegateTool(s *SubAgent) common.Tool {
	// TODO
	return common.Tool{}
}

// Orchestrator 是主 Agent 的 ToolRunner：
//
//	名字为 "delegate_<sub>" 且 sub 存在 -> 解析 {"task"}（失败或 task 为空返回 ErrBadDelegate），调用子 Agent 的 Run；
//	  子 Agent 的结果原样返回，并在 Transcript 追加 "<sub>: <task>"。
//	其它名字 -> 交给 Fallback（Fallback 为 nil 时返回 error "unknown tool: <name>"）。
type Orchestrator struct {
	Subs       map[string]*SubAgent
	Fallback   common.ToolRunner
	Transcript []string
}

func (o *Orchestrator) Call(ctx context.Context, name, argsJSON string) (string, error) {
	// TODO
	return "", nil
}

// Defs 返回所有子 Agent 的 DelegateTool，按子 Agent 名字排序。
func (o *Orchestrator) Defs() []common.Tool {
	// TODO
	return nil
}
