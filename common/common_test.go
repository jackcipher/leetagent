package common

import (
	"context"
	"os"
	"testing"
)

func TestMockServerStitching(t *testing.T) {
	srv := NewMockServer(
		MockReply{Msg: AssistantCalls(Call("call_1", "get_list", `{"page":2,"forum":"科技"}`), Call("call_2", "fetch_one", `{"tid":"123"}`)), ChunkSize: 2},
		MockReply{Msg: AssistantText("你好，世界")},
	)
	defer srv.Close()
	c := NewClient(srv.URL())
	m, err := c.Complete(context.Background(), "m", []Message{{Role: "user", Content: "hi"}}, []Tool{NewFunctionTool("get_list", "d", `{"type":"object"}`)})
	if err != nil {
		t.Fatal(err)
	}
	if len(m.ToolCalls) != 2 || m.ToolCalls[0].Function.Arguments != `{"page":2,"forum":"科技"}` || m.ToolCalls[1].ID != "call_2" {
		t.Fatalf("bad stitch: %+v", m.ToolCalls)
	}
	m, err = c.Complete(context.Background(), "m", nil, nil)
	if err != nil || m.Content != "你好，世界" {
		t.Fatalf("got %q %v", m.Content, err)
	}
	if srv.RequestCount() != 2 || srv.Requests[0].ToolChoice != "auto" {
		t.Fatalf("requests not recorded")
	}
}

// TestLiveLLM 只在 LEETAGENT_LIVE=1 时运行，连真实的 本地 LLM 模型。
func TestLiveLLM(t *testing.T) {
	if os.Getenv("LEETAGENT_LIVE") != "1" {
		t.Skip("set LEETAGENT_LIVE=1 to hit the real LLM")
	}
	base := os.Getenv("LEETAGENT_BASE_URL")
	if base == "" {
		base = "http://localhost:8000/v1"
	}
	model := os.Getenv("LEETAGENT_MODEL")
	if model == "" {
		model = "mlx-community/Qwen2.5-0.5B-Instruct-4bit"
	}
	m, err := NewClient(base).Complete(context.Background(), model, []Message{{Role: "user", Content: "用一句话介绍你自己"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("live reply: %s", m.Content)
}
