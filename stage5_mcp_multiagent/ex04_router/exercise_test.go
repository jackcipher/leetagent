package ex04_router

import (
	"context"
	"errors"
	"testing"

	"leetagent/common"
)

var routes = []Route{
	{Name: "forum", Description: "论坛帖子相关：查帖子、总结讨论", Keywords: []string{"帖子", "论坛", "tid"}},
	{Name: "code", Description: "写代码、解释代码、报错排查", Keywords: []string{"代码", "panic", "golang", "Go "}},
	{Name: "chitchat", Description: "闲聊", Keywords: []string{"你好", "谢谢"}},
}

func TestKeywordRoute(t *testing.T) {
	tests := []struct {
		text, want string
		ok         bool
	}{
		{"帮我总结论坛里关于 GMP 的帖子", "forum", true},
		{"我的 Golang 代码 PANIC 了", "code", true},
		{"你好，论坛的代码在哪", "forum", true}, // forum 1, code 1, chitchat 1 -> 靠前
		{"今天天气怎么样", "", false},
	}
	for _, tc := range tests {
		got, ok := KeywordRoute(routes, tc.text)
		if got != tc.want || ok != tc.ok {
			t.Errorf("%q: got %q %v", tc.text, got, ok)
		}
	}
	dup := []Route{{Name: "a", Keywords: []string{"x", "x", "X"}}, {Name: "b", Keywords: []string{"x", "y"}}}
	if got, _ := KeywordRoute(dup, "x y"); got != "b" {
		t.Fatalf("duplicate keywords count once: got %q", got)
	}
}

func TestClassifyPrompt(t *testing.T) {
	want := "你是一个意图分类器。根据用户输入选择最合适的类别，只输出类别名，不要输出其它内容。\n类别：\n- forum: 论坛帖子相关：查帖子、总结讨论\n- code: 写代码、解释代码、报错排查\n- chitchat: 闲聊\n"
	if got := ClassifyPrompt(routes); got != want {
		t.Fatalf("got:\n%s", got)
	}
}

func TestRouter(t *testing.T) {
	tests := []struct {
		name, text string
		reply      common.Reply
		want, src  string
		llmCalls   int
	}{
		{"keyword wins", "看看 tid 12", common.Reply{}, "forum", "keyword", 0},
		{"llm exact", "天气", common.Reply{Msg: common.AssistantText("  Chitchat。\n")}, "chitchat", "llm", 1},
		{"llm backticks", "天气", common.Reply{Msg: common.AssistantText("`code`")}, "code", "llm", 1},
		{"llm sentence", "天气", common.Reply{Msg: common.AssistantText("我认为属于 forum 类别")}, "forum", "llm", 1},
		{"llm ambiguous", "天气", common.Reply{Msg: common.AssistantText("forum 或 code")}, "chitchat", "default", 1},
		{"llm garbage", "天气", common.Reply{Msg: common.AssistantText("不知道")}, "chitchat", "default", 1},
		{"llm error", "天气", common.Reply{Err: errors.New("timeout")}, "chitchat", "default", 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			llm := &common.ScriptedChatter{Replies: []common.Reply{tc.reply}}
			r := &Router{LLM: llm, Model: "glm", Routes: routes, Default: "chitchat"}
			got, src := r.Route(context.Background(), tc.text)
			if got != tc.want || src != tc.src || llm.CallCount() != tc.llmCalls {
				t.Fatalf("got %q %q calls=%d", got, src, llm.CallCount())
			}
			if tc.llmCalls == 1 {
				c := llm.Calls[0]
				if c.Tools != nil || c.Msgs[0].Content != ClassifyPrompt(routes) || c.Msgs[1].Content != tc.text {
					t.Fatalf("llm call = %+v", c)
				}
			}
		})
	}
}
