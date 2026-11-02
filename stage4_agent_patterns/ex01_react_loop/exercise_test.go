package ex01_react_loop

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"leetagent/common"
)

type fakeTools struct{ calls []string }

func (f *fakeTools) Call(ctx context.Context, name, args string) (string, error) {
	f.calls = append(f.calls, name+args)
	if name == "boom" {
		return "", errors.New("forum down")
	}
	return "result of " + name, nil
}

// 这个测试走真正的 HTTP + SSE：mock 服务器把 tool_call 参数拆成多片，用 common.Client 拼接。
func TestRunOverMockHTTP(t *testing.T) {
	srv := common.NewMockServer(
		common.MockReply{Msg: common.AssistantCalls(
			common.Call("c1", "get_list", `{"page":1}`),
			common.Call("c2", "boom", `{}`),
		), ChunkSize: 2},
		common.MockReply{Msg: common.AssistantText("第一页有 3 个帖子，论坛另一个接口挂了。")},
	)
	defer srv.Close()
	tools := &fakeTools{}
	a := &Agent{LLM: common.NewClient(srv.URL()), Tools: tools, Model: "glm", MaxRounds: 5,
		Defs: []common.Tool{common.NewFunctionTool("get_list", "列表", `{"type":"object"}`)}}
	in := []common.Message{{Role: "system", Content: "sys"}, {Role: "user", Content: "看第一页"}}
	answer, hist, err := a.Run(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	if answer != "第一页有 3 个帖子，论坛另一个接口挂了。" {
		t.Fatalf("answer = %q", answer)
	}
	if !reflect.DeepEqual(tools.calls, []string{`get_list{"page":1}`, "boom{}"}) {
		t.Fatalf("tool calls = %v", tools.calls)
	}
	wantRoles := []string{"system", "user", "assistant", "tool", "tool", "assistant"}
	var roles []string
	for _, m := range hist {
		roles = append(roles, m.Role)
	}
	if !reflect.DeepEqual(roles, wantRoles) {
		t.Fatalf("roles = %v", roles)
	}
	if hist[3].ToolCallID != "c1" || hist[4].ToolCallID != "c2" || hist[4].Content != "error: forum down" {
		t.Fatalf("tool messages wrong: %+v %+v", hist[3], hist[4])
	}
	if len(in) != 2 {
		t.Fatal("input modified")
	}
	// 第二次请求必须带上 assistant.tool_calls 和两条 tool 结果
	second := srv.Requests[1]
	if len(second.Messages) != 5 || len(second.Messages[2].ToolCalls) != 2 || len(second.Tools) != 1 {
		t.Fatalf("second request = %+v", second)
	}
}

func TestRunMaxRounds(t *testing.T) {
	loop := common.AssistantCalls(common.Call("c", "get_list", "{}"))
	fake := common.NewScripted(loop, loop, loop, loop)
	a := &Agent{LLM: fake, Tools: &fakeTools{}, MaxRounds: 3}
	_, _, err := a.Run(context.Background(), []common.Message{{Role: "user", Content: "x"}})
	if !errors.Is(err, ErrMaxRounds) || fake.CallCount() != 3 {
		t.Fatalf("err=%v calls=%d", err, fake.CallCount())
	}
}

func TestRunLLMError(t *testing.T) {
	fake := &common.ScriptedChatter{Replies: []common.Reply{{Err: errors.New("EOF")}}}
	a := &Agent{LLM: fake, Tools: &fakeTools{}, MaxRounds: 3}
	if _, _, err := a.Run(context.Background(), nil); err == nil || err.Error() != "EOF" {
		t.Fatalf("err = %v", err)
	}
}

func TestParsePlan(t *testing.T) {
	out := "好的，计划如下：\n1. 获取第一页列表\n  2、 挑出 Go 相关帖子\n3) 逐个阅读\n10. 写总结\n- 不是步骤\n4.\n"
	want := []string{"获取第一页列表", "挑出 Go 相关帖子", "逐个阅读", "写总结"}
	if got := ParsePlan(out); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %q", got)
	}
}

func TestPlanAndExecute(t *testing.T) {
	fake := common.NewScripted(
		common.AssistantText("1. 列出帖子\n2. 写总结"),
		common.AssistantCalls(common.Call("c", "get_list", "{}")),
		common.AssistantText("拿到列表"),
		common.AssistantText("总结完毕"),
	)
	a := &Agent{LLM: fake, Tools: &fakeTools{}, MaxRounds: 3, Model: "m"}
	res, err := a.PlanAndExecute(context.Background(), "调研 Go 帖子")
	if err != nil || !reflect.DeepEqual(res, []string{"拿到列表", "总结完毕"}) {
		t.Fatalf("res=%v err=%v", res, err)
	}
	if fake.Calls[0].Tools != nil || fake.Calls[0].Msgs[0].Content != PlanPrompt {
		t.Fatal("planning call should use PlanPrompt and no tools")
	}
	if !strings.Contains(fake.Calls[3].Msgs[0].Content, "第 2 步（共 2 步）：写总结") {
		t.Fatalf("step prompt = %q", fake.Calls[3].Msgs[0].Content)
	}
}
