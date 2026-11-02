package m3_research_agent

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"leetagent/capstone/m1_forum_tools/solution"
	"leetagent/common"
)

// 测试用 M1 的参考实现作为工具，这样 M3 不依赖你 M1 是否写完。
func forum(t *testing.T) common.ToolRunner {
	f, err := solution.LoadForum("../m1_forum_tools/testdata/forum.json")
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestSystemPrompt(t *testing.T) {
	if !strings.Contains(SystemPrompt, "[tid]") {
		t.Fatal("SystemPrompt must explain the [tid] citation format")
	}
}

func TestExtractCitations(t *testing.T) {
	tests := []struct {
		in   string
		want []int
	}{
		{"用 pprof [101]，加 ctx [101][104]。", []int{101, 104}},
		{"[ 101 ] [abc] [] 无", nil},
		{"见[103]和[7]", []int{103, 7}},
	}
	for _, tc := range tests {
		got := ExtractCitations(tc.in)
		if len(got) == 0 && len(tc.want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%q → %v want %v", tc.in, got, tc.want)
		}
	}
}

func TestMarkdown(t *testing.T) {
	r := Report{Question: "Q", Answer: "A [1]", Verified: []int{1}, Unverified: []int{9}}
	want := "# 研究：Q\n\nA [1]\n\n## 引用\n- [1] ✅\n- [9] ⚠️ 未验证\n"
	if got := r.Markdown(); got != want {
		t.Fatalf("got %q\nwant %q", got, want)
	}
	if got := (Report{Question: "Q", Answer: "A"}).Markdown(); got != "# 研究：Q\n\nA\n\n## 引用\n（无）\n" {
		t.Fatalf("got %q", got)
	}
}

func TestRun(t *testing.T) {
	llm := common.NewScripted(
		common.AssistantCalls(common.Call("c1", "search_posts", `{"query":"goroutine"}`)),
		common.AssistantCalls(common.Call("c2", "get_post", `{"tid":101}`), common.Call("c3", "get_post", `{"tid":555}`)),
		common.AssistantText("用 pprof 排查 [101]，另见 [999]。"),
	)
	r := &Researcher{LLM: llm, Model: "m", Tools: forum(t), Defs: []common.Tool{common.NewFunctionTool("get_post", "d", `{}`)}}
	rep, err := r.Run(context.Background(), "goroutine 泄漏怎么查？")
	if err != nil {
		t.Fatal(err)
	}
	if rep.Answer != "用 pprof 排查 [101]，另见 [999]。" || rep.ToolCalls != 3 ||
		!reflect.DeepEqual(rep.Verified, []int{101}) || !reflect.DeepEqual(rep.Unverified, []int{999}) {
		t.Fatalf("report = %+v", rep)
	}
	if llm.CallCount() != 3 {
		t.Fatalf("llm calls = %d", llm.CallCount())
	}
	last := llm.Calls[2]
	if last.Msgs[0].Role != "system" || last.Msgs[1].Content != "goroutine 泄漏怎么查？" || len(last.Tools) != 1 {
		t.Fatalf("bad request: %+v", last.Msgs[:2])
	}
	// msgs: sys,user,a1,t1,a2,t2,t3
	if len(last.Msgs) != 7 {
		t.Fatalf("history len = %d", len(last.Msgs))
	}
	if m := last.Msgs[3]; m.Role != "tool" || m.ToolCallID != "c1" || !strings.Contains(m.Content, "105") {
		t.Fatalf("tool msg = %+v", m)
	}
	if m := last.Msgs[6]; m.ToolCallID != "c3" || m.Content != "error: post 555 not found" {
		t.Fatalf("tool error msg = %+v", m)
	}
}

func TestRunLimits(t *testing.T) {
	loop := common.NewScripted(
		common.AssistantCalls(common.Call("a", "get_list", `{}`)),
		common.AssistantCalls(common.Call("b", "get_list", `{}`)),
		common.AssistantText("never reached"),
	)
	r := &Researcher{LLM: loop, Tools: forum(t), MaxRounds: 2}
	if _, err := r.Run(context.Background(), "q"); !errors.Is(err, ErrMaxRounds) {
		t.Fatalf("err = %v, want ErrMaxRounds", err)
	}
	boom := errors.New("upstream down")
	bad := &common.ScriptedChatter{Replies: []common.Reply{{Err: boom}}}
	r = &Researcher{LLM: bad, Tools: forum(t)}
	if _, err := r.Run(context.Background(), "q"); !errors.Is(err, boom) {
		t.Fatalf("err = %v", err)
	}
}
