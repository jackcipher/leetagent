// Package solution 是 ex01_react_loop 的参考答案。
package solution

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"leetagent/common"
)

var ErrMaxRounds = errors.New("agent exceeded max rounds")

type Agent struct {
	LLM       common.Chatter
	Tools     common.ToolRunner
	Defs      []common.Tool
	Model     string
	MaxRounds int
}

func (a *Agent) Run(ctx context.Context, msgs []common.Message) (string, []common.Message, error) {
	history := make([]common.Message, len(msgs), len(msgs)+8)
	copy(history, msgs)
	for round := 1; round <= a.MaxRounds; round++ {
		reply, err := a.LLM.Complete(ctx, a.Model, history, a.Defs)
		if err != nil {
			return "", history, err
		}
		reply.Role = "assistant"
		history = append(history, reply)
		if len(reply.ToolCalls) == 0 {
			return reply.Content, history, nil
		}
		for _, call := range reply.ToolCalls {
			history = append(history, a.runTool(ctx, call))
		}
	}
	return "", history, ErrMaxRounds
}

func (a *Agent) runTool(ctx context.Context, call common.ToolCall) common.Message {
	out, err := a.Tools.Call(ctx, call.Function.Name, call.Function.Arguments)
	if err != nil {
		out = "error: " + err.Error()
	}
	return common.Message{Role: "tool", ToolCallID: call.ID, Content: out}
}

const PlanPrompt = "先不要调用工具。请把用户的目标拆成 3~6 个可执行的步骤，每行一个，格式为“1. 步骤”。只输出步骤列表。"

func parseStep(line string) (string, bool) {
	line = strings.TrimSpace(line)
	i := 0
	for i < len(line) && line[i] >= '0' && line[i] <= '9' {
		i++
	}
	if i == 0 || i >= len(line) {
		return "", false
	}
	rest := line[i:]
	for _, sep := range []string{".", "、", ")", "）"} {
		if strings.HasPrefix(rest, sep) {
			step := strings.TrimFunc(rest[len(sep):], unicode.IsSpace)
			return step, step != ""
		}
	}
	return "", false
}

func ParsePlan(output string) []string {
	var steps []string
	for _, line := range strings.Split(output, "\n") {
		if s, ok := parseStep(line); ok {
			steps = append(steps, s)
		}
	}
	return steps
}

func (a *Agent) PlanAndExecute(ctx context.Context, goal string) ([]string, error) {
	plan, err := a.LLM.Complete(ctx, a.Model, []common.Message{
		{Role: "system", Content: PlanPrompt},
		{Role: "user", Content: goal},
	}, nil)
	if err != nil {
		return nil, err
	}
	steps := ParsePlan(plan.Content)
	var results []string
	for i, step := range steps {
		msgs := []common.Message{
			{Role: "system", Content: fmt.Sprintf("你在执行计划的第 %d 步（共 %d 步）：%s", i+1, len(steps), step)},
			{Role: "user", Content: goal},
		}
		ans, _, err := a.Run(ctx, msgs)
		if err != nil {
			return results, err
		}
		results = append(results, ans)
	}
	return results, nil
}
