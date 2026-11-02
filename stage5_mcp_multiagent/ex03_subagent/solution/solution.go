// Package solution 是 ex03_subagent 的参考答案。
package solution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"leetagent/common"
)

var (
	ErrMaxRounds   = errors.New("subagent exceeded max rounds")
	ErrBadDelegate = errors.New("delegate arguments must be {\"task\": \"...\"} with non-empty task")
)

type SubAgent struct {
	Name        string
	Description string
	System      string
	LLM         common.Chatter
	Model       string
	Tools       common.ToolRunner
	Defs        []common.Tool
	MaxRounds   int
}

func (s *SubAgent) Run(ctx context.Context, task string) (string, error) {
	history := []common.Message{{Role: "system", Content: s.System}, {Role: "user", Content: task}}
	for i := 0; i < s.MaxRounds; i++ {
		reply, err := s.LLM.Complete(ctx, s.Model, history, s.Defs)
		if err != nil {
			return "", err
		}
		reply.Role = "assistant"
		history = append(history, reply)
		if len(reply.ToolCalls) == 0 {
			return reply.Content, nil
		}
		for _, c := range reply.ToolCalls {
			history = append(history, s.runTool(ctx, c))
		}
	}
	return "", ErrMaxRounds
}

func (s *SubAgent) runTool(ctx context.Context, c common.ToolCall) common.Message {
	var out string
	var err error
	if s.Tools == nil {
		err = fmt.Errorf("unknown tool: %s", c.Function.Name)
	} else {
		out, err = s.Tools.Call(ctx, c.Function.Name, c.Function.Arguments)
	}
	if err != nil {
		out = "error: " + err.Error()
	}
	return common.Message{Role: "tool", ToolCallID: c.ID, Content: out}
}

const delegateSchema = `{"type":"object","properties":{"task":{"type":"string","description":"交给子 Agent 的完整任务描述"}},"required":["task"]}`

func DelegateTool(s *SubAgent) common.Tool {
	return common.NewFunctionTool("delegate_"+s.Name, s.Description, delegateSchema)
}

type Orchestrator struct {
	Subs       map[string]*SubAgent
	Fallback   common.ToolRunner
	Transcript []string
}

func (o *Orchestrator) Call(ctx context.Context, name, argsJSON string) (string, error) {
	if sub, ok := o.Subs[strings.TrimPrefix(name, "delegate_")]; ok && strings.HasPrefix(name, "delegate_") {
		var args struct {
			Task string `json:"task"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil || strings.TrimSpace(args.Task) == "" {
			return "", ErrBadDelegate
		}
		o.Transcript = append(o.Transcript, sub.Name+": "+args.Task)
		return sub.Run(ctx, args.Task)
	}
	if o.Fallback == nil {
		return "", fmt.Errorf("unknown tool: %s", name)
	}
	return o.Fallback.Call(ctx, name, argsJSON)
}

func (o *Orchestrator) Defs() []common.Tool {
	names := make([]string, 0, len(o.Subs))
	for n := range o.Subs {
		names = append(names, n)
	}
	sort.Strings(names)
	defs := make([]common.Tool, len(names))
	for i, n := range names {
		defs[i] = DelegateTool(o.Subs[n])
	}
	return defs
}
