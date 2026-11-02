package ex03_subagent

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"leetagent/common"
)

type forumTools struct{}

func (forumTools) Call(ctx context.Context, name, args string) (string, error) {
	if name == "get_post" {
		return `{"tid":7,"title":"GMP 调度"}`, nil
	}
	return "", errors.New("no such tool")
}

func TestSubAgentRun(t *testing.T) {
	llm := common.NewScripted(
		common.AssistantCalls(common.Call("a", "get_post", `{"tid":7}`), common.Call("b", "nope", `{}`)),
		common.AssistantText("帖子 7 讲 GMP 调度。"),
	)
	s := &SubAgent{Name: "researcher", System: "你是研究员", LLM: llm, Model: "glm", Tools: forumTools{}, MaxRounds: 3}
	out, err := s.Run(context.Background(), "读帖子 7")
	if err != nil || out != "帖子 7 讲 GMP 调度。" {
		t.Fatalf("out=%q err=%v", out, err)
	}
	first := llm.Calls[0].Msgs
	if len(first) != 2 || first[0].Content != "你是研究员" || first[1].Content != "读帖子 7" {
		t.Fatalf("subagent must start with fresh context: %+v", first)
	}
	second := llm.Calls[1].Msgs
	if second[3].Content != `{"tid":7,"title":"GMP 调度"}` || second[4].Content != "error: no such tool" {
		t.Fatalf("tool results = %+v", second[3:])
	}

	loop := common.AssistantCalls(common.Call("a", "get_post", "{}"))
	s2 := &SubAgent{LLM: common.NewScripted(loop, loop), MaxRounds: 2}
	if _, err := s2.Run(context.Background(), "x"); !errors.Is(err, ErrMaxRounds) {
		t.Fatalf("err = %v", err)
	}
	s3 := &SubAgent{LLM: common.NewScripted(common.AssistantCalls(common.Call("a", "x", "{}")), common.AssistantText("ok")), MaxRounds: 2}
	if out, err := s3.Run(context.Background(), "x"); err != nil || out != "ok" {
		t.Fatalf("nil Tools must not panic: %q %v", out, err)
	}
}

func TestDelegateTool(t *testing.T) {
	d := DelegateTool(&SubAgent{Name: "writer", Description: "把要点写成中文报告"})
	if d.Type != "function" || d.Function.Name != "delegate_writer" || d.Function.Description != "把要点写成中文报告" {
		t.Fatalf("d = %+v", d)
	}
	var schema struct {
		Required   []string                  `json:"required"`
		Properties map[string]map[string]any `json:"properties"`
	}
	if err := json.Unmarshal(d.Function.Parameters, &schema); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(schema.Required, []string{"task"}) || schema.Properties["task"]["type"] != "string" {
		t.Fatalf("schema = %s", d.Function.Parameters)
	}
}

type recorder struct{ names []string }

func (r *recorder) Call(ctx context.Context, name, args string) (string, error) {
	r.names = append(r.names, name)
	return "fallback:" + name, nil
}

func TestOrchestrator(t *testing.T) {
	writer := &SubAgent{Name: "writer", Description: "写报告", LLM: common.NewScripted(common.AssistantText("报告完成")), MaxRounds: 2}
	researcher := &SubAgent{Name: "researcher", Description: "调研", MaxRounds: 1}
	fb := &recorder{}
	o := &Orchestrator{Subs: map[string]*SubAgent{"writer": writer, "researcher": researcher}, Fallback: fb}
	var _ common.ToolRunner = o

	out, err := o.Call(context.Background(), "delegate_writer", `{"task":"写一段总结"}`)
	if err != nil || out != "报告完成" {
		t.Fatalf("out=%q err=%v", out, err)
	}
	if !reflect.DeepEqual(o.Transcript, []string{"writer: 写一段总结"}) {
		t.Fatalf("transcript = %v", o.Transcript)
	}
	for _, bad := range []string{`{}`, `{"task":"  "}`, `not json`} {
		if _, err := o.Call(context.Background(), "delegate_writer", bad); !errors.Is(err, ErrBadDelegate) {
			t.Errorf("%s: err = %v", bad, err)
		}
	}
	if out, _ := o.Call(context.Background(), "get_list", "{}"); out != "fallback:get_list" {
		t.Fatalf("fallback out = %q", out)
	}
	if out, _ := o.Call(context.Background(), "delegate_ghost", `{"task":"x"}`); out != "fallback:delegate_ghost" {
		t.Fatalf("unknown sub should go to fallback, got %q", out)
	}
	o.Fallback = nil
	if _, err := o.Call(context.Background(), "x", "{}"); err == nil || err.Error() != "unknown tool: x" {
		t.Fatalf("err = %v", err)
	}
	defs := o.Defs()
	if len(defs) != 2 || defs[0].Function.Name != "delegate_researcher" || defs[1].Function.Name != "delegate_writer" {
		t.Fatalf("defs = %+v", defs)
	}
}
