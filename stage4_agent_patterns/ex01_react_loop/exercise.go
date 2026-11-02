//go:build !solution

package ex01_react_loop

import (
	"context"
	"errors"

	"leetagent/common"
)

var ErrMaxRounds = errors.New("agent exceeded max rounds")

// Agent 是带工具的 ReAct 循环。
type Agent struct {
	LLM       common.Chatter
	Tools     common.ToolRunner
	Defs      []common.Tool
	Model     string
	MaxRounds int
}

// Run 执行循环（和 你的 Agent 项目 的 runOneRound 同构）：
//
//	每轮：Complete(ctx, Model, msgs, Defs) 得到 assistant 消息并追加到历史；
//	没有 ToolCalls -> 返回 (该消息的 Content, 完整历史, nil)；
//	有 ToolCalls -> 按顺序执行每个调用，追加 {role:"tool", tool_call_id, content}；
//	  工具出错时 content = "error: " + err.Error()（把错误交给模型自己处理，而不是中断）；
//	超过 MaxRounds 轮仍在调工具 -> 返回 ("", 历史, ErrMaxRounds)。
//
// 模型调用出错直接返回该错误。入参 msgs 不能被修改。
func (a *Agent) Run(ctx context.Context, msgs []common.Message) (string, []common.Message, error) {
	// TODO
	return "", nil, nil
}

// PlanPrompt 让模型先列计划。
const PlanPrompt = "先不要调用工具。请把用户的目标拆成 3~6 个可执行的步骤，每行一个，格式为“1. 步骤”。只输出步骤列表。"

// ParsePlan 从模型输出中提取步骤：识别以 "数字." 或 "数字、" 或 "数字)" 开头的行（允许前导空白），
// 返回去掉编号和首尾空白后的文本；其它行忽略。
func ParsePlan(output string) []string {
	// TODO
	return nil
}

// PlanAndExecute 先用 PlanPrompt 作为 system、goal 作为 user 调一次模型（不带工具）得到计划，
// 然后对每个步骤依次调用 a.Run，消息为 [{system: "你在执行计划的第 i 步（共 n 步）：<step>"}, {user: goal}]，
// i 从 1 开始。返回每一步的最终回答。任一步出错立即返回已完成的结果和错误。计划为空返回 (nil, nil)。
func (a *Agent) PlanAndExecute(ctx context.Context, goal string) ([]string, error) {
	// TODO
	return nil, nil
}
